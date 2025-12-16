package server

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func quickRunTask(ctx *gin.Context) {
	// 设置SSE头部
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")
	tool := ctx.Param("tool")

	var (
		result string
		err    error
	)
	defer func() {
		closeSSE(ctx, result)
	}()

	var printMessage func(string) = newPrintMessage(ctx)

	if tool == "" {
		printMessage("tool is required")
		return
	}

	toolFunc, ok := toolMapping[tool]
	if !ok {
		printMessage(fmt.Sprintf("tool %s not found", tool))
		return
	}

	input := queryAll(ctx)
	bytes, _ := json.Marshal(input)

	result, err = toolFunc(string(bytes), printMessage)
	if err != nil {
		printMessage("")
		printMessage(fmt.Sprintf("failed to run task: %v", err))
		return
	}
	if len(result) > 0 {
		printMessage("task result: " + result)
	}
}

func newPrintMessage(ctx *gin.Context) func(string) {
	return func(msg string) {
		ctx.SSEvent("message", msg)
		ctx.Writer.Flush()
	}
}

func closeSSE(ctx *gin.Context, msg string) {
	ctx.SSEvent("close", msg)
	ctx.Writer.Flush()
}
