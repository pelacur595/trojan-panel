package main

import (
	"github.com/gin-gonic/gin"
	"trojan/router"
)

func main() {
	r := gin.Default()
	router.Router(r)
	_ = r.Run(":8081")
}
