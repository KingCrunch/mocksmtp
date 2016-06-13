package mail

import (
	"time"
	"net/textproto"
	"github.com/satori/go.uuid"
	"strings"
	"mime/multipart"
	"bitbucket.org/chrj/smtpd"
	"mime"
	"io"
	"io/ioutil"
	"bufio"
	"log"
)

type Mail struct {
	Id uuid.UUID

	ReceivedAt time.Time
	Sender string
	Recipients []string
	Header textproto.MIMEHeader
	Data []byte

	Multipart bool
	Parts []MailPart
}

type MailPart struct {
	Header textproto.MIMEHeader
	Disposition string
	DispositionParams map[string]string
	Data []byte
}


func NewMail (peer smtpd.Peer, env smtpd.Envelope) Mail {
	reader := bufio.NewReader(strings.NewReader(string(env.Data)))
	tp := textproto.NewReader(reader)

	mimeHeader, err := tp.ReadMIMEHeader()
	check(err)

	mediaType, params, err := mime.ParseMediaType(mimeHeader.Get("Content-Type"))
	check(err)

	m := &Mail{
		Id: uuid.NewV4(),
		ReceivedAt: time.Now(),
		Sender: string(env.Sender),
		Recipients: env.Recipients,
		Header: mimeHeader,
		Data: env.Data,
	}



	if strings.HasPrefix(mediaType, "multipart/") {
		m.Multipart = true
		m.Parts = make([]MailPart, 0, 0)

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

			m.Parts = append(m.Parts, *part)
		}
	}

	return *m
}


func check (err error) {
	if err != nil {
		log.Fatal(err)
	}
}
