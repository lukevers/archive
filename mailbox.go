package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sg3des/eml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Mailbox struct {
	Dirs   []string
	Emails []*Email
}

func NewMailbox(m *Mailbox) *Mailbox {
	m.index()
	go m.watch()
	return m
}

func (m *Mailbox) watch() {
	watcher, err := fsnotify.NewWatcher()
	defer watcher.Close()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println(event)
				err := m.index()
				if err != nil {
					panic(err)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	for _, dir := range m.Dirs {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	<-done
}

func (m *Mailbox) index() error {
	m.Emails = nil
	emails := make(map[string]*Email)

	err := filepath.Walk("./emails", func(path string, f os.FileInfo, err error) error {
		// Stop if it's not an EML file.
		if !strings.HasSuffix(path, ".eml") {
			return nil
		}

		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		message, err := eml.Parse(bytes)
		if err != nil {
			return err
		}

		email := &Email{Message: &message}
		m.Emails = append(m.Emails, email)
		emails[message.MessageId] = email

		return nil
	})

	if err != nil {
		return err
	}

	// Sort the slice of emails by date
	sort.Slice(m.Emails, func(i, j int) bool {
		return m.Emails[i].Message.Date.After(m.Emails[j].Message.Date)
	})

	for _, email := range m.Emails {
		// Loop over each reply that this email has and link it to the previous
		// email.
		for _, id := range email.Message.InReply {
			if e, exists := emails[id]; exists {
				e.Next = email
				email.Prev = e
			}
		}
	}

	for _, email := range m.Emails {
		log.Println(email.Message.Subject)
		log.Println(email)
	}

	return nil
}
