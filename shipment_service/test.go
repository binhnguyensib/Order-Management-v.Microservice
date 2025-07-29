package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	formattedDate := now.Format("20060102")
	result := "price_update_" + formattedDate + ".xlsx"
	fmt.Println("Generated filename:", result)
}
