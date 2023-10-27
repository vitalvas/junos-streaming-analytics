package main

import (
	"context"
	"log"
	"time"

	"github.com/vitalvas/gokit/xcmd"
	"github.com/vitalvas/junos-streaming-analytics/internal/core"
	"golang.org/x/sync/errgroup"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	group, ctx := errgroup.WithContext(context.Background())

	conf, err := core.ParseCollectorConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	app, err := core.NewCore(ctx, conf)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < conf.Workers; i++ {
		group.Go(app.RunJTIProcess)
	}

	group.Go(func() error {
		return xcmd.PeriodicRun(ctx, app.SendAllOutputs, time.Second)
	})

	if conf.JTI != nil {
		group.Go(app.RunJTI)
	}

	group.Go(func() error {
		err := xcmd.WaitInterrupted(ctx)
		app.Shutdown(ctx)
		return err
	})

	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
}
