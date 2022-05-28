package abstract

import (
	"fmt"

	"github.com/draymonders/bingo/utils"
)

type Obj interface {
	GetId() int64
}

type ICompany interface {
	GetCompanyName() string
}

type IPerson interface {
	GetPersonName() string
}

type Subject struct {
	Id          int64  `utils:"id"`
	PersonName  string `utils:"person_name"`
	CompanyName string `utils:"company_name"`
}

func (s *Subject) GetId() int64 {
	return s.Id
}

func (s *Subject) GetPersonName() string {
	return s.PersonName
}

// 根据传入的obj进行打印
func print(obj Obj) {
	if obj == nil {
		fmt.Println("obj is null")
		return
	}
	fmt.Printf("format obj: %v\n", utils.ObjToJson(obj))

	if er, ok := obj.(ICompany); ok {
		companyName := er.GetCompanyName()
		fmt.Println("companyName: ", companyName)
	}
	if er, ok := obj.(IPerson); ok {
		personName := er.GetPersonName()
		fmt.Println("personName: ", personName)
	}
}

// interface转换
func RunItfTrans() {
	var obj Obj
	print(obj)
	fmt.Printf("====\n")

	obj = &Subject{
		Id:          1,
		PersonName:  "22",
		CompanyName: "333",
	}
	print(obj)
	fmt.Printf("====\n")
}
