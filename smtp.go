package main

import (
	"bitbucket.org/chrj/smtpd"
	"github.com/satori/go.uuid"
	"log"
	"bufio"
	"strings"
	"net/textproto"
	"mime"
	"time"
	"mime/multipart"
	"io"
	"io/ioutil"
)

func RunSmtpServer(bind string, ch chan<- struct {smtpd.Peer; smtpd.Envelope}) {
	server := &smtpd.Server{
		Handler: func(peer smtpd.Peer, env smtpd.Envelope) error {
			log.Printf("New Mail from %q to %q", env.Sender, env.Recipients)
			ch<- struct{smtpd.Peer; smtpd.Envelope}{peer, env}

			return nil
		},
	}

	err := server.ListenAndServe(bind)
	check(err)
}


func handEnvelope (env smtpd.Envelope) {
	reader := bufio.NewReader(strings.NewReader(string(env.Data)))
	tp := textproto.NewReader(reader)

	mimeHeader, err := tp.ReadMIMEHeader()
	check(err)

	mediaType, params, err := mime.ParseMediaType(mimeHeader.Get("Content-Type"))
	check(err)

	mail := &Mail{
		Id: uuid.NewV4(),
		ReceivedAt: time.Now(),
		Sender: string(env.Sender),
		Recipients: env.Recipients,
		Header: mimeHeader,
		Data: env.Data,
	}



	if strings.HasPrefix(mediaType, "multipart/") {
		mail.Multipart = true
		mail.Parts = make([]MailPart, 0, 0)

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

			part := &MailPart{
				Header: p.Header,
				Disposition: disp,
				DispositionParams: dispParams,
				Data: slurp,
			}

			mail.Parts = append(mail.Parts, *part)
		}
	}

	mailBucket[mail.Id] = *mail
}
