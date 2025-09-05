package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"redis-clone/config"
	"redis-clone/core"
)

func readCommand(conn net.Conn) (*core.Cmd, error) {
	//TODO: can only read a buffer of 4096 need to write a repeated read till EOF/delimiter logic
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	tokens, err := core.DecodeArrayString(buffer[:n])
	if err != nil {
		return nil, err
	}
	Cmd := &core.Cmd{
		Cmd:  tokens[0],
		Args: tokens[1:],
	}

	return Cmd, err
}

func writeCommand(conn net.Conn, val []byte) error {
	_, err := conn.Write(val)

	return err
}
func writeErrorCommand(conn net.Conn, err error) error {
	_, errWrite := conn.Write(core.EncodeError(err))
	return errWrite
}

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
