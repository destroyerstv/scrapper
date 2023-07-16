package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	sitechecker "scrapper/internal/checker"
	"scrapper/internal/service"
	datastorage "scrapper/internal/storage"
	"syscall"
	"time"
)

func main() {
	ds, err := datastorage.NewDataStorage("../sites.txt")

	if err != nil {
		fmt.Println("datastorage error")
	}

	client := http.Client{Timeout: 60 * time.Second}
	ctx, cancel := context.WithCancel(context.Background())

	sc := sitechecker.NewSiteChecker(ctx, client, ds)
	go func() {
		sc.Run()
	}()

	exit := make(chan os.Signal, 1)

	s := service.NewService(ds)
	go func() {
		err := s.Run("8080")
		if err != nil {
			fmt.Println(err)
			exit <- syscall.SIGTERM
		}
	}()

	// graceful shutdown
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	cancel()
}
