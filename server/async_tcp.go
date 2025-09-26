package server

import (
	"fmt"
	"log"
	"net"
	"redis-clone/config"
	"redis-clone/core"
	"syscall"
)

/*
eventloop -->continously listenes if we are getting any events that we can listen to.
Kqueue--> emits the events if there is a i/o

- Create Non blocking listener
- create a kqueue
- register the listener with the kqueue
- Evenet loop just listenes if we are getting any kqueue events
*/
type Socket struct {
	socketFd int
}

func Listen(ip string, port int) (*Socket, error) {
	s := &Socket{}
	listenFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)

	if err != nil {

		return nil, err
	}

	s.socketFd = listenFd
	sockAddr := &syscall.SockaddrInet4{Port: port}
	copy(sockAddr.Addr[:], net.ParseIP(ip))
	err = syscall.Bind(s.socketFd, sockAddr)
	if err != nil {
		return nil, fmt.Errorf("error binding address with socket %v", err)
	}
	err = syscall.Listen(s.socketFd, syscall.SOMAXCONN)
	if err != nil {
		return nil, fmt.Errorf("failed to listen")
	}
	return s, err
}

func registerKqueue(socket *Socket, kfd int) error {
	event := syscall.Kevent_t{
		Ident:  uint64(socket.socketFd),
		Filter: syscall.EVFILT_READ, //Monitor for read events
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE,
		Fflags: 0,
		Data:   0,
		Udata:  nil,
	}
	_, err := syscall.Kevent(kfd, []syscall.Kevent_t{event}, nil, nil)
	return err
}
func eventLoop(socket *Socket, kfd int) {
	events := make([]syscall.Kevent_t, 2500)
	for {
		n, err := syscall.Kevent(kfd, nil, events, nil)
		if err != nil {
			log.Printf("Kevent error: %v", err)
			continue
		}
		for i := 0; i < n; i++ {
			fd := int(events[i].Ident)
			//if my server is ready for I/O rather than client (Means new client ready to connect)
			if fd == socket.socketFd {
				//accept the incoming connection from client
				clientSocket := &Socket{}
				clientFd, _, err := syscall.Accept(socket.socketFd)
				if err != nil {
					log.Printf("client connection error: %v", err)
				}
				clientSocket.socketFd = clientFd
				syscall.SetNonblock(clientFd, true)
				fmt.Printf("new client connected ", clientFd)
				//register the new client socket/fd in the kqueue
				registerKqueue(clientSocket, kfd)
			} else {
				//if the i/o is from any other socket/fd rather than our server socket/fd
				//just perform read and write
				c := &FDcomm{Fd: fd}
				cmd, err := readCommand(c)
				if err != nil {
					syscall.Close(int(fd))
					continue
				}
				output, err := core.Eval(cmd)
				var errWrite error
				if err != nil {
					errWrite = writeErrorCommand(c, err)
				} else {
					errWrite = writeCommand(c, output)
				}
				if errWrite != nil {
					log.Println("err:% write", err)
				}
			}
		}
	}
}
func RunAsyncTcpServer() error {
	socket, err := Listen(config.Host, config.Port)
	if err != nil {
		return err
	}
	fmt.Println("I am here")
	defer syscall.Close(socket.socketFd)
	kfd, err := syscall.Kqueue()
	if err != nil {
		return err
	}
	err = registerKqueue(socket, kfd)
	if err != nil {
		return err
	}
	fmt.Println("Server started")
	eventLoop(socket, kfd)
	return nil
}
