package main

import "github.com/PaulYakow/metrics-track/internal/app"

func main() {
	client := app.NewClient()
	client.Run()
}
