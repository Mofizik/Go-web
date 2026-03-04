package main

import (
    "order/internal/app"
	"context"
)

func main() {
	ctx := context.Background()
    a, err := app.New(ctx)
	if err != nil {
		return
	}
	if err := a.Run(); err != nil {
		panic(err)
	}
	a.Run()
}