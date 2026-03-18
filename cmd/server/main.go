package main

import (
	"context"
	"order/internal/app"
)

func main() {
	ctx := context.Background()
	a, err := app.New(ctx)
	if err != nil {
		panic(err)
	}
	if err := a.Run(); err != nil {
		panic(err)
	}
}
