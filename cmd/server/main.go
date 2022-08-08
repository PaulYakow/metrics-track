package main

import (
	"github.com/PaulYakow/metrics-track/internal/app"
)

func main() {
	server := app.NewServer()
	server.Run()
}
