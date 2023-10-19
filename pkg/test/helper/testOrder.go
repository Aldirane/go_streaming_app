package helper

import (
	"encoding/json"
	"fmt"
	"go_app/pkg/order"
)

// Print generated fake order
func PrintFakeData() {
	order := order.GenerateFakeData()
	data, err := json.MarshalIndent(order, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", data)

}
