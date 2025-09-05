package main

import (
	"flag"
	"fmt"
	"log"
	"redis-clone/config"
	"redis-clone/server"
)

func init() {
	// Set log flags to include file and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for redis server")
	flag.StringVar(&config.Port, "port", "6379", "host for redis server")
	flag.Parse()
}

func main() {
	setupFlags()
	fmt.Println("Starting db....")
	server.RunSyncTcpServer()
	// out, _ := core.Decode([]byte("*-1\r\n"))
	// fmt.Println(out)

}

//
