package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	inputMapping map[string]string = make(map[string]string)
	locker       sync.Mutex        = sync.Mutex{}
)

func createInput(ctx *gin.Context) {
	inputKey := fmt.Sprintf("%d", time.Now().UnixNano())
	bytes, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to get raw data: %v", err)})
		return
	}
	locker.Lock()
	defer locker.Unlock()
	inputMapping[inputKey] = string(bytes)
	ctx.JSON(http.StatusOK, gin.H{"input_key": inputKey})
}

func getInput(ctx *gin.Context) string {
	query := queryAll(ctx)
	bytes, _ := json.Marshal(query)
	inputKey, ok := query["input_key"]
	if !ok {
		return string(bytes)
	}
	if value, ok := inputMapping[inputKey]; ok {
		return value
	}

	return string(bytes)
}

func queryAll(ctx *gin.Context) map[string]string {
	query := ctx.Request.URL.Query()
	result := make(map[string]string)
	for key := range query {
		result[key] = query.Get(key)
	}
	return result
}
