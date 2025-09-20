package server

import (
	"io"
	"redis-clone/core"
	"syscall"
)

type FDcomm struct {
	Fd int
}

func (f *FDcomm) Read(b []byte) (int, error) {
	return syscall.Read(f.Fd, b)
}
func (f *FDcomm) Write(b []byte) (int, error) {
	return syscall.Write(f.Fd, b)
}
func readCommand(c io.ReadWriter) (*core.Cmd, error) {
	//TODO: can only read a buffer of 4096 need to write a repeated read till EOF/delimiter logic
	buffer := make([]byte, 4096)
	n, err := c.Read(buffer)
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

func writeCommand(c io.ReadWriter, val []byte) error {
	_, err := c.Write(val)

	return err
}
func writeErrorCommand(c io.ReadWriter, err error) error {
	_, errWrite := c.Write(core.EncodeError(err))
	return errWrite
}
