package stock

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const productNum = 100

var db *gorm.DB

/*
库存扣减，需要依赖 mysql耶

1. 单纯用goroutine测试下表的tps
2. 使用本地队列削峰试一下
*/

func Init() {
	dsn := "root:root@tcp(10.248.162.48:3306)/scenes?charset=utf8mb4&parseTime=True&loc=Local"
	var err error

	cfg := &gorm.Config{}

	db, err = gorm.Open(mysql.Open(dsn), cfg)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(productNum)
	sqlDB.SetMaxOpenConns(productNum * 100)
	sqlDB.SetConnMaxLifetime(time.Minute)

	dbErr := db.Transaction(func(tx *gorm.DB) error {
		var curProductNum int64
		db.Table("t_stock").Count(&curProductNum)
		if curProductNum > 0 {
			err = db.Table("t_stock").Where("id > 0").Delete("*").Error
			if err != nil {
				return err
			}
		}

		stocks := make([]*stockModel, 0, productNum)
		for i := 1; i <= productNum; i++ {
			stocks = append(stocks, &stockModel{
				ProductId: productId(i),
				StockNum:  0,
			})
		}
		err = db.Table("t_stock").CreateInBatches(stocks, productNum).Error
		if err != nil {
			return err
		}
		return nil
	})
	if dbErr != nil {
		panic(dbErr)
	}

}

// IncrementV1 用goroutine测试下库存表的tps
// @param n 表示有多少个人去竞争行锁
func IncrementV1(n int) {
	Init()

	// 测算下timeout的时间内能更新多少数据
	timeout := 10
	goRoutineNums := n * productNum
	chs := make(chan struct{}, goRoutineNums)
	wg := sync.WaitGroup{}

	// 开100个协程去update
	for i := 0; i < goRoutineNums; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			vId := id%productNum + 1
			mId := id / productNum
			pId := productId(vId)

			jobPrefix := fmt.Sprintf("job ProductId %v(%v) ", pId, mId)

			for {
				select {
				case <-chs:
					fmt.Printf("%v done\n", jobPrefix)
					return
				default:
					dbErr := db.Transaction(func(tx *gorm.DB) error {
						db.Table("t_stock").Where("product_id = ?", pId).Updates(map[string]interface{}{"stock_num": gorm.Expr("stock_num + ?", 1)})
						return nil
					})
					if dbErr != nil {
						//fmt.Printf("%v err %v\n", jobPrefix, dbErr)
						// 稍微休息下
						time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
					}
				}
			}
		}(i)
	}

	timeOutCh := time.After(time.Duration(timeout) * time.Second)
	select {
	case <-timeOutCh:
		for i := 1; i <= goRoutineNums; i++ {
			chs <- struct{}{}
		}
		fmt.Printf("!!! timeout, so close all sub jobs\n")
	}

	wg.Wait()
	// 都结束了的话，查询下总cnt

	calcCount(timeout)

	// 结果
	// => stock_nums 5031 used 10s
	// => stock_nums 4359 used 10s
	// => stock_nums 4464 used 10s
	// tps 大概在 400 - 500 op/s
}

// calcCount 统计timeout时间内执行的库存加减数
func calcCount(timeout int) {
	cnts := make([]int64, 0)
	if err := db.Table("t_stock").Where("id > ?", 0).Select("stock_num").Find(&cnts).Error; err != nil {
		fmt.Printf("count all stock_num err: %v\n", cnts)
		os.Exit(1)
	}
	cnt := int64(0)
	for _, v := range cnts {
		cnt += v
	}

	fmt.Printf("stock_nums %v used %vs\n", cnt, timeout)
}

func IncrementV2(ctx context.Context) {
}

type stockModel struct {
	Id        int64 `json:"id"`
	ProductId int64 `json:"product_id"`
	StockNum  int64 `json:"stock_num"`
}

func productId(i int) int64 {
	return int64(1000 + i)
}
