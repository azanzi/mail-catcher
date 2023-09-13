package main

import (
	"flag"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strings"
)

// config holds the cmd parameters
type config struct {
	smtpPort int
	webPort  int
}

// application represents the state of the app.
// It contains information related to the server and data.
type application struct {
	config config
	// inbox stores all the emails received during runtime.
	inbox []*Email
	// templates is the cache containing all html templates preloaded in memory.
	templates map[string]*template.Template
}

// listen listens for incoming TCP connections and handles them asynchronously.
func (app *application) listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.config.smtpPort))
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
	email := NewEmail()

	// Welcome message
	c.writeLine(220, "Mail Catcher")

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
		if isHeader(cmd) {
			email.parseHeader(cmd)
		} else {
			email.addHeader(cmd)
		}
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

	fmt.Println("Mail successfully received")
}

func main() {
	var cfg config

	flag.IntVar(&cfg.smtpPort, "smtp", 1025, "SMTP server port")
	flag.IntVar(&cfg.webPort, "web", 9999, "Web interface port")
	flag.Parse()

	app := &application{
		config:    cfg,
		templates: make(map[string]*template.Template),
	}
	app.loadTemplates()

	go http.ListenAndServe(fmt.Sprintf(":%d", app.config.webPort), app.routes())

	fmt.Printf("SMTP: localhost:%d\n", app.config.smtpPort)
	fmt.Printf("Web UI: http://localhost:%d\n", app.config.webPort)
	err := app.listen()
	if err != nil {
		panic(err)
	}
}
