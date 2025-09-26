package core

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
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
func evalSET(cmd *Cmd) ([]byte, error) {
	if len(cmd.Args) < 2 {
		return nil, errors.New("ERR Wrong number of arguments for ping command")
	}
	key, value := cmd.Args[0], cmd.Args[1]
	val := &Value{
		value:  value,
		expiry: -1,
	}
	for i := 2; i < len(cmd.Args); i++ {
		switch cmd.Args[i] {
		case "EX", "ex":
			if i+1 >= len(cmd.Args) {
				return nil, errors.New("ERR syntax error")
			}
			expiry, err := strconv.ParseInt(cmd.Args[i+1], 10, 64)
			if err != nil {
				return nil, errors.New("ERR value is not an integer or out of range")
			}
			val.expiry = time.Now().Unix() + expiry

			i++
		}
	}
	Put(key, val)
	return Encode("OK", true), nil
}
func evalGET(cmd *Cmd) ([]byte, error) {
	if len(cmd.Args) < 1 {
		return nil, errors.New("ERR Wrong number of arguments for ping command")
	}
	key := cmd.Args[0]
	val := Get(key)
	fmt.Println("expiry:", val.expiry)
	if val == nil || val.expiry <= time.Now().Unix() {
		return RESPNIL, nil
	}
	return Encode(val.value, false), nil
}
func evalTTL(cmd *Cmd) ([]byte, error) {
	if len(cmd.Args) < 1 {
		return nil, errors.New("ERR Wrong number of arguments for ping command")
	}

	key := cmd.Args[0]
	val := Get(key)
	if val.expiry == -1 {
		return Encode(val.value, false), nil
	}
	TTL := val.expiry - time.Now().Unix()

	return Encode(TTL, false), nil

}
func Eval(cmd *Cmd) ([]byte, error) {
	switch strings.ToUpper(cmd.Cmd) {
	case "PING":
		return evalPing(cmd)
	case "GET":
		return evalGET(cmd)
	case "SET":
		return evalSET(cmd)
	case "TTL":
		return evalTTL(cmd)
	}
	return Encode("#{cmd.Cmd} Unknown", true), nil
}
