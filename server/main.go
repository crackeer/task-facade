package server

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/joho/godotenv/autoload"
)

var (
	toolMapping map[string]func(string, func(string)) (string, error) = make(map[string]func(string, func(string)) (string, error))
)

func enableCORS(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}

func Run(tools map[string]func(string, func(string)) (string, error), port string) {
	for key, tempFunc := range tools {
		toolMapping[key] = tempFunc
	}
	// 创建 Gin 实例
	router := gin.Default()
	router.Use(enableCORS, gin.Logger(), gin.Recovery())
	router.GET("/run/:tool", runTask)
	router.GET("/quick_run/:tool", quickRunTask)

	if len(port) < 1 {
		if value, ok := os.LookupEnv("PORT"); ok {
			port = value
		} else {
			port = "8080"
		}
	}
	// 启动服务器
	router.Run(":" + port)
}
