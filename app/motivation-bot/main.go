package main

import (
	"motivation-bot/di"
)

func main() {
	container := di.BuildContainer()
	Start(container)
}
