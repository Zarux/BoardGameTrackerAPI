package main

import (
	"flag"
	"github.com/Zarux/BGServer/internal/app/BGServer"
	"log"
	"os"
)

var (
	port       = flag.String("port", "8080", "http service address")
	configFile = flag.String("config", "configs/conf.json", "path to config file")
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.SetOutput(os.Stdout)
	flag.Parse()
	BGServer.Run(*port, *configFile)
}
