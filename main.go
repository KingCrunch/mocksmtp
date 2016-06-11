package main

//go:generate go-bindata-assetfs static/... templates/...

import (
	"github.com/satori/go.uuid"
	"log"
	"flag"
	"fmt"
	"runtime"
	"path/filepath"
	"os"
	"os/signal"
	"github.com/KingCrunch/visualsmtp/store"
	"github.com/KingCrunch/visualsmtp/mail"
)

const Name string = "visualsmtp"
const Version string = "0.1.0"
const Help string = `
This tool provides a SMTP-Server and a HTTP-server.
`

var (
	mailBucket map[uuid.UUID]mail.Mail = make(map[uuid.UUID] mail.Mail)
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

	s := store.NewInMemoryStore()

	fmt.Println("Start HTTP-Server listening on "+options.HttpBind)
	go RunHttpServer(options.HttpBind, s)

	fmt.Println("Start SMTP-server listening on "+options.SmtpBind)
	go RunSmtpServer(options.SmtpBind, s)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		select {
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
