package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"redis-clone/config"
	"redis-clone/core"
)

func RunSyncTcpServer() {
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	fmt.Printf("Starting Server at %s", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		for {
			Command, err := readCommand(conn)
			// conn.Read()
			if err != nil {
				conn.Close()
				if err == io.EOF {
					break
				}
				log.Println("client_disconnnected")
			}
			fmt.Println("OUTPUT", Command)
			//Eval the command
			//write the output
			output, err := core.Eval(Command)
			log.Println(output)
			var errWrite error
			if err != nil {
				errWrite = writeErrorCommand(conn, err)
			} else {
				errWrite = writeCommand(conn, output)
			}
			if errWrite != nil {
				log.Println("err:% write", err)
			}

		}

	}
}
