package main

import (
	"flag"
	"github.com/Zarux/BGServer/internal/app/server"
	"log"
	"os"
)

var port = flag.String("port", "8080", "http service address")

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.SetOutput(os.Stdout)
	flag.Parse()
	server.Run(*port)
}
