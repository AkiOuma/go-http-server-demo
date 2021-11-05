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
	mux := http.NewServeMux()
	mux.HandleFunc("/", home.HomePage)
	mux.HandleFunc("/stop", stopper.StopServer)

	// config server object and start
	app := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}

	// monitor if server need to stop because of receving interupt signal
	go func() {
		<-ctx.Done()
		if err := app.Shutdown(ctx); err != nil {
			serverLog.Printf("HTTP server(%d) Shutdown: \n\t%v", port, err)
		}
	}()

	// simulate if server stop by user calling stop api
	go func() {
		<-c
		serverLog.Printf("HTTP server(%d) Shutdown by api", port)
		if err := app.Shutdown(ctx); err != nil {
			serverLog.Printf("HTTP server(%d) Shutdown: \n\t%v", port, err)
		}
	}()
	log.Printf("starting http server at %d", port)
	return app.ListenAndServe()
}
