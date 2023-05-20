package main

import (
	"go.uber.org/dig"
	"motivation-bot/di"
)

func Start(container *dig.Container) {
	err := container.Invoke(func(app *di.App) {
		app.Client.Run()
	})
	if err != nil {
		panic(err)
	}
}
