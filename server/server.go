package server

import (
	"context"
	"http-server-demo/server/controller"
	"log"
	"net/http"
	"os"
	"strconv"
)

func ServerRun(ctx context.Context, port int) error {
	serverLog := log.New(os.Stdout, "http server    : ", log.Ldate|log.Ltime)
	c := make(chan struct{})
	// initialize controllers
	home := controller.NewHome()
	stopper := controller.NewStopper(c)

	// register handler
	http.HandleFunc("/", home.HomePage)
	http.HandleFunc("/stop", stopper.StopServer)

	// config server object and start
	app := &http.Server{
		Addr: ":" + strconv.Itoa(port),
	}

	// monitor if server need to stop because of receving interupt signal
	go func() {
		<-ctx.Done()
		if err := app.Shutdown(ctx); err != nil {
			serverLog.Printf("HTTP server Shutdown: \n\t%v", err)
		}
	}()

	// simulate if server stop by user calling stop api
	go func() {
		<-c
		serverLog.Println("HTTP server Shutdown by api")
		if err := app.Shutdown(ctx); err != nil {
			serverLog.Printf("HTTP server Shutdown: \n\t%v", err)
		}
	}()

	return app.ListenAndServe()
}
