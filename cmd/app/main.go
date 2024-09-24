package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/dostonshernazarov/mini-twitter/internal/app"
	configpkg "github.com/dostonshernazarov/mini-twitter/internal/pkg/config"
)

func main() {
	// config
	config := configpkg.Load()

	// app
	app, err := app.NewApp(*config)
	if err != nil {
		log.Fatal(err)
	}

	// run application
	go func() {
		app.Logger.Info("Listen: ", zap.String("address", config.Server.Host+config.Server.Port))
		if err := app.Run(); err != nil {

			app.Logger.Error("app run", zap.Error(err))
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// app stops
	app.Logger.Info("api gateway service stops")
	app.Stop()
}
