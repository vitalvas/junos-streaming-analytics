package core

import (
	"context"
	"fmt"

	"github.com/vitalvas/junos-streaming-analytics/internal/output"
	"github.com/vitalvas/junos-streaming-analytics/internal/output/console"
)

const (
	udpReadBufferSize = 16384
	jtiBacklog        = 10000
)

type App struct {
	ctx     context.Context
	config  *CollectorConfig
	outputs map[string]output.Output

	jtiRouter   map[string][]string // map of JTI input to output names
	jtiMessages chan jtiMessage

	shutdownCh chan struct{}
}

func NewCore(ctx context.Context, config *CollectorConfig) (*App, error) {
	core := &App{
		ctx:         ctx,
		config:      config,
		jtiRouter:   make(map[string][]string),
		jtiMessages: make(chan jtiMessage, jtiBacklog),
		outputs:     make(map[string]output.Output),
		shutdownCh:  make(chan struct{}),
	}

	for name, outputConfig := range config.Outputs {
		switch outputConfig.Type {
		case "console":
			output, err := console.NewOutput(outputConfig)
			if err != nil {
				return nil, err
			}

			core.outputs[name] = output

		// case "prometheus":
		// 	output, err := prometheus.NewOutput(outputConfig)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// 	core.outputs[name] = output

		default:
			return nil, fmt.Errorf("unknown output type: %s", outputConfig.Type)
		}
	}

	if core.outputs == nil {
		return nil, fmt.Errorf("no outputs defined")
	}

	return core, nil
}

func (app *App) Shutdown(ctx context.Context) {
	close(app.shutdownCh)

	app.SendAllOutputs(ctx)
}

func (app *App) SendAllOutputs(_ context.Context) error {
	for _, output := range app.outputs {
		if err := output.Send(); err != nil {
			return err
		}
	}

	return nil
}
