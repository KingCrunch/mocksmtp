package main

//go:generate go-bindata-assetfs static/... templates/...

import (
	"bitbucket.org/chrj/smtpd"
	"github.com/satori/go.uuid"
	"log"
	"strings"
	"net/textproto"
	"bufio"
	"mime"
	"mime/multipart"
	"io"
	"io/ioutil"
	"time"
	"flag"
	"fmt"
	"runtime"
	"path/filepath"
	"os"
)

const Name string = "visualsmtp"
const Version string = "0.1.0"
const Help string = `
This tool provides a SMTP-Server and a HTTP-server.
`

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

var (
	mailBucket map[uuid.UUID]Mail = make(map[uuid.UUID] Mail)
	options struct {
		HttpBind string
		SmtpBind string
		Version bool
		Help bool
	}
)

func init() {
	flag.StringVar(&options.HttpBind, "http-bind", ":12080", "The IP and port to bind the HTTP-server to: [<ip>]:port")
	flag.StringVar(&options.SmtpBind, "smtp-bind", ":12025", "The IP and port to bind the SMTP-server to: [<ip>]:port")
	flag.BoolVar(&options.Help, "help", false, "Help")
	flag.BoolVar(&options.Help, "h", false, "See --help")
	flag.BoolVar(&options.Version, "version", false, "Show Version")

}

func main() {
	flag.Parse()

	if (options.Help) {
		file, err := filepath.Abs(os.Args[0])
		check(err)
		fmt.Println(Name+"-"+Version+"-"+runtime.GOOS+"-"+runtime.GOARCH)
		fmt.Println(file+" [-help|-h] [-version] [-http-bind=[<ip>]:<port>] [-smtp-bind=[<ip>]:<port>]")
		fmt.Println(Help)
		flag.PrintDefaults()
		return
	}

	if (options.Version) {
		fmt.Println(Name+"-"+Version+"-"+runtime.GOOS+"-"+runtime.GOARCH)
		return
	}

	go RunHttpServer(options.HttpBind)

	server := &smtpd.Server{
		Handler: func(peer smtpd.Peer, env smtpd.Envelope) error {
			log.Printf("New Mail from %q to %q", env.Sender, env.Recipients)
			go handEnvelope(env)

			return nil
		},
	}
	err := server.ListenAndServe(options.SmtpBind)
	check(err)
}

func check (err error) {
	if err != nil {
		log.Fatal(err)
	}
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
