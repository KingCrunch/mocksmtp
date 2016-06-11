package main

import (
	"bitbucket.org/chrj/smtpd"
	"log"
	"bufio"
	"strings"
	"net/textproto"
	"mime"
	"time"
	"mime/multipart"
	"io"
	"io/ioutil"
	"github.com/KingCrunch/visualsmtp/mail"
	"github.com/KingCrunch/visualsmtp/store"
)

func RunSmtpServer(bind string, s store.Store) {
	server := &smtpd.Server{
		Handler: func(peer smtpd.Peer, env smtpd.Envelope) error {
			log.Printf("New Mail from %q to %q", env.Sender, env.Recipients)
			m := handEnvelope(env)
			_, err := s.Push(m)
			check(err)

			return nil
		},
	}

	err := server.ListenAndServe(bind)
	check(err)
}


func handEnvelope (env smtpd.Envelope) mail.Mail {
	reader := bufio.NewReader(strings.NewReader(string(env.Data)))
	tp := textproto.NewReader(reader)

	mimeHeader, err := tp.ReadMIMEHeader()
	check(err)

	mediaType, params, err := mime.ParseMediaType(mimeHeader.Get("Content-Type"))
	check(err)

	m := &mail.Mail{
		ReceivedAt: time.Now(),
		Sender: string(env.Sender),
		Recipients: env.Recipients,
		Header: mimeHeader,
		Data: env.Data,
	}



	if strings.HasPrefix(mediaType, "multipart/") {
		m.Multipart = true
		m.Parts = make([]mail.MailPart, 0, 0)

		mr := multipart.NewReader(strings.NewReader(string(env.Data)), params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			check(err)
			slurp, err := ioutil.ReadAll(p)
			check(err)

			disp, dispParams, err  := mime.ParseMediaType(p.Header.Get("Content-Disposition"))
			check(err)

			part := &mail.MailPart{
				Header: p.Header,
				Disposition: disp,
				DispositionParams: dispParams,
				Data: slurp,
			}

			m.Parts = append(m.Parts, *part)
		}
	}

	return *m
}
