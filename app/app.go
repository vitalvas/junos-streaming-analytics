package app

import (
	"context"
	"os/signal"
	"runtime"
	"syscall"
)

const Backlog = 10000

type App struct {
	message chan []byte
}

func NewApp() *App {
	return &App{
		message: make(chan []byte, Backlog),
	}
}

func (app *App) Run() {
	go app.receiver()

	for i := 0; i <= runtime.NumCPU(); i++ {
		go app.process()
	}
}

func Execute() {
	app := NewApp()

	go app.Run()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
}
