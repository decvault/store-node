package main

import (
	"context"
	"log"
)

func main() {
	app := App()
	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("failed to start app: %+v", err)
	}

	<-app.Done()
	if err := app.Stop(context.Background()); err != nil {
		log.Fatalf("error occured on app stop: %+v", err)
	}
}
