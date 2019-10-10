package main

import (
	"cjvirtucio87/distributed-todo-go/pkg/config"
)

func main() {
	loader := config.NewViperLoader(
		"app",
		"yaml",
	)

	loader.Load()
}
