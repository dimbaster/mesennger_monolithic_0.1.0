package main

import "Server/internal/app"

func main() {
	application := app.New()
	application.Run()
}
