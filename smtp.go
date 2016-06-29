package main

import (
	"bytes"
	"log"

	"bitbucket.org/chrj/smtpd"
	"github.com/KingCrunch/mocksmtp/store"
	"github.com/veqryn/go-email/email"
)

func RunSmtpServer(bind string, s store.Store) {
	server := &smtpd.Server{
		Handler: func(peer smtpd.Peer, env smtpd.Envelope) error {
			log.Printf("New Mail from %q to %q", env.Sender, env.Recipients)
			m, _ := email.ParseMessage(bytes.NewReader(env.Data))
			err := s.Push(store.NewItem(m, env.Sender, env.Recipients))

			return err
		},
	}

	err := server.ListenAndServe(bind)

	if err != nil {
		log.Fatal(err)
	}
}
