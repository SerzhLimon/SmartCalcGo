package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"APG2_SmartCalc/app/internal/server"
	storage "APG2_SmartCalc/app/pkg/redis"
)

func main() {

	storage := storage.NewRedisCache()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := storage.GetClient().Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Redis is running on :6379")

	server := server.NewServer(storage)
	r := gin.Default()

	r.LoadHTMLFiles("./frontend/index.html")

	r.GET("/", server.ServeHTML)
	r.POST("/push", server.PushNumber)
	r.POST("/pushx", server.PushX)

	r.POST("/pushoperator", server.PushOperator)
	r.POST("/pushfunction", server.PushFunction)

	r.POST("/pushpoint", server.PushPoint)

	r.POST("/pushopenbracket", server.PushOpenBracket)
	r.POST("/pushclosebracket", server.PushCloseBracket)

	r.POST("/pushsign", server.PushChangeSign)
	r.DELETE("/pushclear", server.PushClear)
	r.DELETE("/pushdeletenode", server.PushDeleteBackNode)

	r.POST("/calc", server.Calc)

	r.POST("/graph", server.MakePlots)
	r.GET("/history", server.GetHistory)

	r.Run(":8000")
}
