package itf

import (
	"fmt"

	"github.com/draymonders/bingo/json"
)

type Obj interface {
	GetId() int64
}

type CompanyNameEr interface {
	GetCompanyName() string
}

type PersonNameEr interface {
	GetPersonName() string
}

type Subject struct {
	Id          int64  `json:"id"`
	PersonName  string `json:"person_name"`
	CompanyName string `json:"company_name"`
}

func (s *Subject) GetId() int64 {
	return s.Id
}

func (s *Subject) GetPersonName() string {
	return s.PersonName
}

// 根据传入的obj进行打印
func format(obj Obj) {
	if obj == nil {
		fmt.Println("obj is null")
		return
	}
	fmt.Printf("format obj: %v\n", json.ObjToJson(obj))

	if er, ok := obj.(CompanyNameEr); ok {
		companyName := er.GetCompanyName()
		fmt.Println("companyName: ", companyName)
	}
	if er, ok := obj.(PersonNameEr); ok {
		personName := er.GetPersonName()
		fmt.Println("personName: ", personName)
	}
}

// interface转换
func RunItfTrans() {
	var obj Obj
	format(obj)
	fmt.Printf("====\n")

	obj = &Subject{
		Id:          1,
		PersonName:  "22",
		CompanyName: "333",
	}
	format(obj)
	fmt.Printf("====\n")
}
