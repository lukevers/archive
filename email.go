package main

import (
	"github.com/sg3des/eml"
)

type Email struct {
	Next *Email
	Prev *Email

	Message *eml.Message
}
