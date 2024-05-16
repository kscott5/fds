package main

import (
	"fmt"
	"github.com/kscott5/fds/orders/services"
)

func main() {
	fmt.Println("orders main")
	services.Create(nil,nil)
}