package main

import (
	"fmt"
)

const configPath = "config/config.yaml"

func main() {
	fmt.Println("Start")
	conf := config.Load(configPath)

	fmt.Println(conf)

	fmt.Println("End")
}
