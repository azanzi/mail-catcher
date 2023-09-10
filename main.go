package main

import (
	"fmt"
	"net"
	"strings"
)

type application struct {
	inbox []Email
	addr  string
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
		go app.handleConnection(NewConn(c))
	}
}

// handleConnection receives a mail from connection c.
// It implements a minimal version of the SMTP protocol.
func (app *application) handleConnection(c *Conn) {
	defer c.Close()
	var email Email

	// Welcome message
	c.writeLine(220, fmt.Sprintf("%s Mail Catcher", app.addr))

	// EHLO/HELO
	cmd, _ := c.readLine()
	c.writeLine(250, "OK")

	// MAIL FROM
	cmd, _ = c.readLine()
	email.SMTP.From = cmd[11 : len(cmd)-1]
	c.writeLine(250, "OK")

	// RCPT TO
	var rcpt []string
	for {
		cmd, _ = c.readLine()
		if !strings.HasPrefix(cmd, "RCPT TO") {
			break
		}
		rcpt = append(rcpt, cmd[9:len(cmd)-1])
		c.writeLine(250, "OK")
	}
	email.SMTP.To = rcpt

	c.writeLine(354, "Start mail input; end with <CRLF>.<CRLF>")

	// DATA
	var body string
	cmd, _ = c.readLine()

	// Read headers
	for cmd != "" {
		// TODO: parse header
		cmd, _ = c.readLine()
	}

	// Read body
	cmd, _ = c.readLine()
	for cmd != "." {
		body += "\n" + cmd
		cmd, _ = c.readLine()
	}

	c.writeLine(250, "OK")
	email.Body = body

	// QUIT
	c.readLine()
	c.writeLine(221, "Service closing transmission channel")

	// Add newly received email to inbox
	app.inbox = append(app.inbox, email)
}

func main() {

}
