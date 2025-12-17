package server

import (
	"github.com/gin-gonic/gin"
)

type Query struct {
	Input    string `json:"input"`
	Output   string `json:"output"`
	Callback string `json:"callback"`
}

func queruQuery(ctx *gin.Context) (Query, error) {
	var query Query
	query.Input = ctx.Query("input")
	query.Output = ctx.Query("output")
	query.Callback = ctx.Query("callback")
	return query, nil
}

func queryAll(ctx *gin.Context) map[string]string {
	query := ctx.Request.URL.Query()
	result := make(map[string]string)
	for key := range query {
		result[key] = query.Get(key)
	}
	return result
}
