package main

import (
	"net"
)

type application struct {
	emails []Email
	addr   string
}

func (app *application) listen() error {
	l, err := net.Listen("tcp", app.addr)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		// TODO: close connection
		if err != nil {
			// TODO: handle error
			continue
		}
		go handleConnection(NewConn(c))
	}
}

func handleConnection(c *Conn) {
	defer c.Close()
	// TODO: handle connection
}

func main() {

}
