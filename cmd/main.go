package main

import (
	"cjvirtucio87/distributed-todo-go/pkg/config"
	"cjvirtucio87/distributed-todo-go/pkg/manager"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	m := manager.NewHttpManager(
		config.NewViperLoader(
			"app",
			"yaml",
		),
	)

	sig := make(chan os.Signal, 1)

	signal.Notify(
		sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	m.Start()

	<-sig
	m.Stop()
}
