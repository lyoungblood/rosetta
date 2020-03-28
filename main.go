/*
 * Rosetta
 *
 * A standard for blockchain interaction
 *
 * API version: 1.2.3
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/celo-org/rosetta/api"
	"github.com/celo-org/rosetta/celo/client"
	"github.com/celo-org/rosetta/internal/config"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	done := make(chan struct{})

	// Read Configuration Variables
	config.SetupDatadir()
	config.ReadConfig()

	log.Info("Initializing Rosetta...", "chainId", config.Chain.ChainId, "epochSize", config.Chain.EpochSize)

	rpcClient, err := rpc.Dial(config.Node.Uri)
	if err != nil {
		log.Crit("Can't connect to node", "err", err)
	}

	celoClient := client.NewCeloClient(rpcClient)

	go api.StartHttpServer(celoClient, done)

	gotExitSignal := HandleSignals()

	<-gotExitSignal
	close(done)
}

func HandleSignals() <-chan bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Press CTRL-C to stop the process")
	go func() {
		sig := <-sigs
		log.Info("Got Signal, Shutting down...", "signal", sig)
		done <- true
	}()

	return done

}
