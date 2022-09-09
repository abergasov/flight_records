package main

import (
	"flag"
	"flight_records/internal/logger"
	"flight_records/internal/routes"
	"flight_records/internal/service/tracker"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var appPort = flag.String("port", "8080", "App listen port")

func main() {
	flag.Parse()
	log, err := logger.NewAppLogger()
	if err != nil {
		println(fmt.Sprintf("error init logger: %s", err))
		return
	}
	trackerService := tracker.InitRouteTracker()
	app := routes.InitAppRouter(*appPort, trackerService)
	go func() {
		log.Info("starting service", zap.String("port", *appPort))
		if err = app.Run(); err != nil {
			log.Fatal("error start service", err)
		}
	}()

	// register app shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // This blocks the main thread until an interrupt is received
	if err = app.Shutdown(); err != nil {
		log.Fatal("error shutdown service", err)
	}
	log.Info("service was successful shutdown")
}
