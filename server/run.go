package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func runTask(ctx *gin.Context) {
	tool := ctx.Param("tool")
	query, _ := queruQuery(ctx)

	toolFunc, ok := toolMapping[tool]
	if !ok {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("tool %s not found", tool))
		return
	}

	go func(queryData Query, tempFunc func(string, func(string)) (string, error)) {
		callback(queryData.Callback, "running", "task is running")
		result, err := executeTask(queryData, tempFunc)
		if err != nil {
			callback(queryData.Callback, "failed", err.Error())
			return
		}
		callback(queryData.Callback, "success", result)
	}(query, toolFunc)

	ctx.JSON(http.StatusOK, map[string]string{"message": "running"})
}

func executeTask(queryData Query, tempFunc func(string, func(string)) (string, error)) (string, error) {
	input, err := getInputData(queryData.Input)
	fmt.Printf("input: %s", input)
	if err != nil {
		return "", err
	}
	outputMessage := newOutputMessage(queryData.Output)
	return tempFunc(input, outputMessage)
}

func newOutputMessage(outputURL string) func(string) {
	return func(msg string) {
		_, err := http.Post(outputURL, "text/plain", strings.NewReader(msg))
		if err != nil {
			fmt.Printf("failed to send message to %s: %v", outputURL, err)
		}
	}
}

func getInputData(inputURL string) (string, error) {
	httpResp, err := http.Get(inputURL)
	if err != nil {
		return "", err
	}
	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func callback(callbackURL string, status, message string) error {
	if callbackURL == "" {
		return nil
	}
	postBody := map[string]string{
		"status":  status,
		"message": message,
	}
	bytes, err := json.Marshal(postBody)
	if err != nil {
		return err
	}

	_, err = http.Post(callbackURL, "application/json", strings.NewReader(string(bytes)))
	return err
}
