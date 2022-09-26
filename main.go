package main

import "github.com/sisu-network/cardano-proxy/core"

func main() {
	s := core.NewServer()
	s.Run()
}
