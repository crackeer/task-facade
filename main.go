package main

import (
	"task-facade/server"
)

func main() {
	toolMapping := make(map[string]func(string, func(string)) (string, error))
	server.Run(toolMapping, "8080")
}
