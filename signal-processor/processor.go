package signalprocessor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func ReceiveSignal(ctx context.Context, signals []os.Signal) error {
	// register signal for notify
	signalLogger := log.New(os.Stdout, "signal receiver: ", log.Ldate|log.Ltime)
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)

	// waiting signal or context done
	select {
	case sig := <-c:
		fmt.Println()
		signalLogger.Printf("receive signal: \n\t%v", sig)
		if sig == os.Interrupt {
			return errors.New(sig.String())
		}
		return nil
	case <-ctx.Done():
		err := ctx.Err()
		signalLogger.Printf("stop receiving signal: \n\t%v", err)
		close(c)
		return errors.New(err.Error())
	}
}
