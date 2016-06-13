package main

import (
	"bitbucket.org/chrj/smtpd"
	"log"
	"github.com/KingCrunch/visualsmtp/store"
	"github.com/KingCrunch/visualsmtp/mail"
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

