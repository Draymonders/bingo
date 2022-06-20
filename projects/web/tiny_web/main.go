package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	serverV2()
}

func serverV0() {
	http.HandleFunc("/user/123", ModifyUserInfo)
	http.HandleFunc("/user/get", GetUserInfo)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// serverV1 map控制路由
func serverV1() {
	server := NewServerWithHandler(NewHandlerBaseOnMap(), MetricBuilder)
	server.Route("GET", "/user", GetUserInfoV1)

	if err := server.Start(":8080"); err != nil {
		panic(err)
	}
}

// serverV2 路由树控制路由
func serverV2() {
	server := NewServerWithHandler(NewHandlerBasedOnTree(), MetricBuilder)
	server.Route("GET", "/user", GetUserInfoV1)

	if err := server.Start(":8080"); err != nil {
		panic(err)
	}
}

func ModifyUserInfo(w http.ResponseWriter, r *http.Request) {

	bytes, _ := ioutil.ReadAll(r.Body)
	u := &UserInfo{}
	_ = json.Unmarshal(bytes, u)

	w.WriteHeader(http.StatusOK)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {

	res := &UserInfo{
		Id:   123,
		Name: "Draymonder",
	}
	bytes, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
}

func GetUserInfoV1(ctx *Context) {

	res := &UserInfo{
		Id:   234,
		Name: "Draymonder",
	}
	bytes, _ := json.Marshal(res)
	time.Sleep(10 * time.Millisecond)
	ctx.W.WriteHeader(http.StatusOK)
	_, _ = ctx.W.Write(bytes)
}

type UserInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
