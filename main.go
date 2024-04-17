package main

import (
	"sync"

	"github.com/PretendoNetwork/mario-kart-7/ctgp7"
	"github.com/PretendoNetwork/mario-kart-7/nex"
)

var wg sync.WaitGroup

func main() {
	wg.Add(3)

	// TODO - Add gRPC server
	go nex.StartAuthenticationServer()
	go nex.StartSecureServer()
	go ctgp7.StartHttpServer()

	wg.Wait()
}
