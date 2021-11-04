package main

import (
	"context"
	"http-server-demo/server"
	signalprocessor "http-server-demo/signal-processor"
	"log"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	log.SetPrefix("main           : ")

	group := new(errgroup.Group)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start http server
	group.Go(func() error {
		err := server.ServerRun(ctx, 8080)
		if err != nil {
			cancel()
		}
		return err
	})

	// start linux signal mornitor
	group.Go(func() error {
		err := signalprocessor.ReceiveSignal(ctx, []os.Signal{
			os.Interrupt,
		})
		if err != nil {
			cancel()
		}
		return err
	})

	if err := group.Wait(); err != nil {
		log.Printf("Exit Reason: \n\t%v\n", err)
	}
}
