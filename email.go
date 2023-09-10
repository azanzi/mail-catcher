package main

import "strings"

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
