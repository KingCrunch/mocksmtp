package main

//go:generate go-bindata-assetfs static/... templates/...

import (
	"bitbucket.org/chrj/smtpd"
	"github.com/satori/go.uuid"
	"log"
	"net/textproto"
	"time"
	"flag"
	"fmt"
	"runtime"
	"path/filepath"
	"os"
	"os/signal"
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

	fmt.Println("Start HTTP-Server listening on "+options.HttpBind)
	go RunHttpServer(options.HttpBind)

	fmt.Println("Start SMTP-server listening on "+options.SmtpBind)
	ch := make(chan struct {smtpd.Peer; smtpd.Envelope}, 10)
	go RunSmtpServer(options.SmtpBind, ch)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		select {
		case rc := <-ch:
			go handEnvelope(rc.Envelope)
		case s := <-c:
			if (s == os.Interrupt) {
				fmt.Println("\nExit. Bye!")
				return
			}
		}
	}
}

func check (err error) {
	if err != nil {
		log.Fatal(err)
	}
}
