package utils

import (
	"fmt"
	"testing"
)

func Test_Id_Unmarshal(t *testing.T) {
	str := "{\"poi_id_v2\":6980984069106436127}"
	data := map[string]interface{}{}
	err := JsonToObj(str, &data)
	if err != nil {
		fmt.Println("jsonToObj err")
		return
	}
	fmt.Println(data)
	for _, v := range data {
		fmt.Printf("%v", v)
	}
	//	map[poi_id_v2:6.980984069106436e+18]
	//	6980984069106436096.000000
}

func Test_Uint64_To_Float64(t *testing.T) {
	// 这种情况下会出现不一致
	v := float64(uint64(6980984069106436127))
	fmt.Println(v)
	fmt.Printf("%f\n", v)
	// 6.980984069106436e+18
	// 6980984069106436096.000000
}
