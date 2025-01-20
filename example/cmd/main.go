package main

import (
	"github.com/woyow/setupper/example/internal/app"
)

func main() {
	a := app.NewApp()

	if err := a.Run(); err != nil {
		panic(err)
	}
}
