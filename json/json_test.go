package json

import (
	"fmt"
	"testing"
)

func Test_Id_Unmarshal(t *testing.T) {
	str := "{\"poi_id_v2\":6980984069106436127}"
	var data map[string]interface{}
	err := JsonToObj(str, &data)
	if err != nil {
		fmt.Println("json err")
		return
	}
	fmt.Println(data)
	for _, v := range data {
		fmt.Printf("%f", v)
	}
	//	map[poi_id_v2:6.980984069106436e+18]
	//	6980984069106436096.000000
}
