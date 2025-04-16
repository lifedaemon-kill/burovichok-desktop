package main

import (
	"context"
	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/config"
	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/logger/z"
	"log"
	"os/signal"
	"syscall"
)

const configPath = "config/config.yaml"

func main() {
	ctx := context.Background()
	err := bootstrap(ctx)
	if err != nil {
		log.Fatalf("main bootstrap failed: %v", err)
	}
}

func bootstrap(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conf := config.Load(configPath)

	if err := z.InitLogger(conf.Logger); err != nil {
		log.Fatalf("init logger error: %s", err)
	}
	defer z.Log.Sync()
	z.Log.Info("init logger and config success")

	/*
		burApp := app.New()
		burWindow := burApp.NewWindow("burovichok")
		text := widget.NewLabel("hello, burovichok!")
		burWindow.SetContent(text)
		burWindow.Resize(fyne.NewSize(800, 400))
		burWindow.ShowAndRun()
	*/
	<-ctx.Done()
	z.Log.Info("context done, shutting down...")

	//	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	//	defer shutdownCancel()

	//TODO завершить sqlite
	z.Log.Info("shutting down completed")
	return nil
}
