package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"http-server-demo/server/controller"
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
	// simulate if server stop by user calling stop api
	go func() {
		select {
		case <-ctx.Done():
			if err := app.Shutdown(ctx); err != nil {
				serverLog.Printf("HTTP server(%d) Shutdown: \n\t%v", port, err)
			}
			return
		case <-c:
			serverLog.Printf("HTTP server(%d) Shutdown by api", port)
			if err := app.Shutdown(ctx); err != nil {
				serverLog.Printf("HTTP server(%d) Shutdown: \n\t%v", port, err)
			}
			return
		}
	}()

	log.Printf("starting http server at %d", port)
	return app.ListenAndServe()
}
