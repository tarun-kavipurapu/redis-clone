package core

import (
	"errors"
	"log"
)

func evalPing(cmd *Cmd) ([]byte, error) {
	log.Println("evalPing", cmd)

	switch len(cmd.Args) {
	case 0:
		// log.Println("evalPing", cmd)
		return Encode("PONG", true), nil
	case 1:
		return Encode(cmd.Args[0], false), nil
	default:
		return nil, errors.New("ERR Wrong number of arguments for ping command")
	}

}
func Eval(cmd *Cmd) ([]byte, error) {
	switch cmd.Cmd {
	case "PING":
		return evalPing(cmd)
	}
	return Encode("#{cmd.Cmd} Unknown", true), nil
}
