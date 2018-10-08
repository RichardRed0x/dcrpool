package main

import (
	"dnldd/dcrpool/worker"
	"fmt"
	"os"
	"os/signal"
	"runtime"
)

func main() {
	// Use all processor cores.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Listen for interrupt signals.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Load configuration and parse command line. This also initializes logging
	// and configures it accordingly.
	cfg, _, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if logRotator != nil {
			logRotator.Close()
		}
	}()

	log.Infof("Version: %s", version())
	log.Infof("Runtime: Go version %s", runtime.Version())
	log.Infof("Home dir: %s", cfg.HomeDir)
	log.Infof("Started dcrpclient.")

	// Initialize the client.
	pc, err := newClient(cfg)
	if err != nil {
		log.Errorf("Failed to create client: %v", err)
		return
	}

	go pc.processMessages()
	for {
		select {
		case <-interrupt:
			pc.cancel()
			if cfg.MinerType == worker.CPU {
				log.Info("CPU miner done.")
			}
			return
		}
	}
}
