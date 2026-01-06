package main

import (
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	app := App()
	if err := app.Start(ctx); err != nil {
		log.Fatalf("failed to start app: %+v", err)
	}

	<-app.Done()
	if err := app.Stop(ctx); err != nil {
		log.Fatalf("error occured on app stop: %+v", err)
	}
}
