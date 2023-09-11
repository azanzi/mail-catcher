package main

import (
	"regexp"
	"strings"
)

// headerRegex matches only 'Subject', 'From' and 'To' fields.
var headerRegex = regexp.MustCompile("^(From|To|Subject|Date):\\s*(.*)$")

type Address string

// Name returns the name of the Address.
func (a Address) Name() string {
	parts := strings.SplitN(string(a), "<", 2)
	return strings.TrimSpace(parts[0])
}

// Email returns the email of the Address.
func (a Address) Email() string {
	parts := strings.SplitN(string(a), "<", 2)
	return strings.TrimRight(parts[1], ">")
}

// Email represents a received email.
// An Email includes both information obtained from the smtp protocol
// and extracted from the email headers inside the DATA transaction.
type Email struct {
	SMTP struct {
		From string
		To   []string
	}
	Subject string
	From    Address
	To      Address
	Body    string
	Headers map[string]string
}

// parseHeader takes a header string and extracts 'Subject', 'From' and 'To' fields.
func (e *Email) parseHeader(s string) {
	parts := strings.SplitN(s, ": ", 2)
	switch parts[0] {
	case "Subject":
		e.Subject = parts[1]
	case "From":
		e.From = Address(parts[1])
	case "To":
		e.To = Address(parts[1])
	}
}

// isHeader checks if the string is in the SMTP header format.
func isHeader(s string) bool {
	return headerRegex.MatchString(s)
}

// addHeader adds a generic header to the Headers.
func (e *Email) addHeader(s string) {
	parts := strings.SplitN(s, ": ", 2)
	e.Headers[parts[0]] = parts[1]
}
