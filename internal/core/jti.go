package core

import (
	"fmt"
	"log"
	"net"
	"time"
)

type jtiMessage struct {
	Source    net.Addr
	Instance  string
	Data      []byte
	Timestamp time.Time
}

func (app *App) RunJTI() error {
	for name, conf := range app.config.JTI {
		log.Println("starting jti collector", "name", name, "addr", conf.Addr)

		go func(name string, conf CollectorJTIConfig) {
			if err := app.runJTIInstance(name, conf); err != nil {
				log.Fatal(err)
			}
		}(name, conf)
	}

	return nil
}

func (app *App) runJTIInstance(instance string, conf CollectorJTIConfig) error {
	if conf.Output == nil {
		app.jtiRouter[instance] = []string{"default"}
	} else {
		app.jtiRouter[instance] = conf.Output
	}

	for _, output := range app.jtiRouter[instance] {
		if _, ok := app.outputs[output]; !ok {
			return fmt.Errorf("unknown output: %s", output)
		}
	}

	udpAddr, err := net.ResolveUDPAddr("udp", conf.Addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	defer conn.Close()

	for {
		select {
		case <-app.shutdownCh:
			return nil

		default:
			buffer := make([]byte, udpReadBufferSize)

			n, err := conn.Read(buffer)
			if err != nil || n == 0 {
				if err != nil {
					log.Printf("error read packet: %s", err.Error())
				}

				continue
			}

			app.jtiMessages <- jtiMessage{
				Source:    conn.RemoteAddr(),
				Instance:  instance,
				Data:      buffer[0:n],
				Timestamp: time.Now(),
			}
		}
	}
}

func (app *App) addMetricToOutput(instance string, name string, labels map[string]string, value float64, timestamp int64) error {
	outputs, ok := app.jtiRouter[instance]
	if !ok {
		return fmt.Errorf("no output route for instance: %s", instance)
	}

	for _, row := range outputs {
		if err := app.outputs[row].AddMetric(name, labels, value, timestamp); err != nil {
			return err
		}
	}

	return nil
}
