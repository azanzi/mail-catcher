package main

import (
	"flag"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strings"
	"time"
)

// config holds the cmd parameters
type config struct {
	smtpAddr string
	webAddr  string
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
	l, err := net.Listen("tcp", app.config.smtpAddr)
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

	// Get received time
	email.ReceivedAt = time.Now()

	// Add newly received email to inbox
	app.inbox = append([]*Email{email}, app.inbox...)

	fmt.Println("Mail successfully received")
}

func main() {
	var cfg config

	flag.StringVar(&cfg.smtpAddr, "smtp", ":1025", "SMTP server address")
	flag.StringVar(&cfg.webAddr, "web", ":9999", "Web interface address")
	flag.Parse()

	templateCache, err := newTemplateCache()
	if err != nil {
		panic(err)
	}

	app := &application{
		config:    cfg,
		templates: templateCache,
	}

	srv := &http.Server{
		Addr:         app.config.webAddr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go srv.ListenAndServe()

	fmt.Printf("SMTP: %s\n", app.config.smtpAddr)
	fmt.Printf("Web UI: %s\n", app.config.webAddr)
	err = app.listen()
	if err != nil {
		panic(err)
	}
}
