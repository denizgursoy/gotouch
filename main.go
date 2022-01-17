/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"fmt"

	"github.com/denizgursoy/gotouch/cmd"
	"github.com/denizgursoy/gotouch/common"
)

func main() {
	cmd.Execute()

	config := common.GetDefaultConfig()
	fmt.Println(config)
}
