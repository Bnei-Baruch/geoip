package main

import (
	"flag"
)

var (
	// http port
	port    = flag.String("p", "8081", "-p=8081")
)

func main() {
	flag.Parse()
	a := App{}
	a.Initialize()
	a.Run(":" + *port)
}
