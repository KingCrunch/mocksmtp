package main

import (
	"log"

	"bitbucket.org/chrj/smtpd"
	"github.com/KingCrunch/mocksmtp/store"
	"github.com/KingCrunch/mocksmtp/mail"
)

func RunSmtpServer(bind string, s store.Store) {
	server := &smtpd.Server{
		Handler: func(peer smtpd.Peer, env smtpd.Envelope) error {
			log.Printf("New Mail from %q to %q", env.Sender, env.Recipients)
			m := mail.NewMail(peer, env)
			err := s.Push(m)
			check(err)

			return nil
		},
	}

	err := server.ListenAndServe(bind)
	check(err)
}

