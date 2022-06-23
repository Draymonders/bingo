package stock

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const productNum = 10

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

	// 开goRoutineNums个协程去update
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
						// 稍微休息下
						gap()
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

func IncrementV2(n int) {
	Init()

	// 测算下timeout的时间内能更新多少数据
	timeout := 10
	producerNum, consumerNum := n*productNum, productNum
	allNum := producerNum + consumerNum
	stopChs := make(chan struct{}, allNum)
	wg := sync.WaitGroup{}

	bufferSize := 1000

	queues := make([]chan struct{}, productNum+1)
	for i := 0; i <= productNum; i++ {
		queues[i] = make(chan struct{}, bufferSize)
	}

	for i := 1; i <= consumerNum; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			vId := id%productNum + 1
			pId := productId(vId)
			for {
				select {
				case <-stopChs:
					fmt.Printf("consume productId %v done\n", pId)
					return
				default:
					consume(pId, queues[vId])
				}
			}
		}(i)
	}
	// 开 producerNum 个协程去 produce msg
	for i := 0; i < producerNum; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			vId := id%productNum + 1
			pId := productId(vId)

			jobPrefix := fmt.Sprintf("producer ProductId %v ", pId)

			for {
				select {
				case <-stopChs:
					fmt.Printf("%v done\n", jobPrefix)
					return
				default:
					produce(queues[vId])
				}
			}
		}(i)
	}

	// timeouts 运算
	go func() {
		timeOutCh := time.After(time.Duration(timeout) * time.Second)
		select {
		case <-timeOutCh:
			for i := 1; i <= allNum; i++ {
				stopChs <- struct{}{}
			}
			fmt.Printf("!!! timeout, so close all sub jobs\n")
		}
	}()

	wg.Wait()
	// 都结束了的话，查询下总cnt

	calcCount(timeout)

	// => stock_nums 21891 used 10s
	// => stock_nums 21918 used 10s
	// => stock_nums 21795 used 10s
}

// produce 生产
func produce(queue chan struct{}) {
	queue <- struct{}{}
	gap()
}

// consume 消费
func consume(pId int64, queue <-chan struct{}) {
	sum := 0
	timeoutCh := time.After(time.Duration(100) * time.Millisecond)
collect:
	for {
		select {
		case <-timeoutCh:
			break collect
		case <-queue:
			sum++
		}
	}

	dbErr := db.Transaction(func(tx *gorm.DB) error {
		db.Table("t_stock").Where("product_id = ?", pId).Updates(map[string]interface{}{"stock_num": gorm.Expr("stock_num + ?", sum)})
		return nil
	})
	if dbErr != nil {
		gap()
	}
}

type stockModel struct {
	Id        int64 `json:"id"`
	ProductId int64 `json:"product_id"`
	StockNum  int64 `json:"stock_num"`
}

func productId(i int) int64 {
	return int64(1000 + i)
}

func gap() {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
}
