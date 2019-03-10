package main

import (
	"flag"
	"github.com/Zarux/BGServer/internal/app/BGServer"
	"log"
)
var (
	port = flag.String("port", "8080", "http service address")
	configFile = flag.String("config", "configs/conf.json", "path to config file")
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	flag.Parse()
	BGServer.Run(*port, *configFile)
}
