package main

import (
	"context"
	"log"
	"os"

	"http-server-demo/server"
	signalprocessor "http-server-demo/signal-processor"

	"golang.org/x/sync/errgroup"
)

func main() {
	log.SetPrefix("main           : ")

	group := new(errgroup.Group)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start http server from port 8080 to 8082
	for i := 8080; i < 8083; i++ {
		port := i
		group.Go(func() error {
			err := server.ServerRun(ctx, port)
			if err != nil {
				cancel()
			}
			return err
		})
	}

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
