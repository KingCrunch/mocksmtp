package main

//go:generate go-bindata-assetfs static/... templates/...

import (
	"flag"
	"fmt"
	"github.com/KingCrunch/mocksmtp/store"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"time"
)

const Name string = "mocksmtp"
const Version string = "0.1.0"
const Help string = `
This tool provides a SMTP-Server and a HTTP-server.
`

var (
	options struct {
		HttpBind  string
		SmtpBind  string
		Retention time.Duration
		Version   bool
		Help      bool
	}
)

func init() {
	flag.StringVar(&options.HttpBind, "http-bind", ":12080", "The IP and port to bind the HTTP-server to: [<ip>]:port")
	flag.StringVar(&options.SmtpBind, "smtp-bind", ":12025", "The IP and port to bind the SMTP-server to: [<ip>]:port")
	i, _ := time.ParseDuration("10m")
	flag.DurationVar(&options.Retention, "retention-time", i, "Retention time")
	flag.BoolVar(&options.Help, "help", false, "Help")
	flag.BoolVar(&options.Help, "h", false, "See --help")
	flag.BoolVar(&options.Version, "version", false, "Show Version")
}

func main() {
	flag.Parse()

	if options.Help {
		file, err := filepath.Abs(os.Args[0])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(Name + "-" + Version + "-" + runtime.GOOS + "-" + runtime.GOARCH)
		fmt.Println(file + " [-help|-h] [-version] [-http-bind=[<ip>]:<port>] [-smtp-bind=[<ip>]:<port>]")
		fmt.Println(Help)
		flag.PrintDefaults()

		return
	}

	if options.Version {
		fmt.Println(Name + "-" + Version + "-" + runtime.GOOS + "-" + runtime.GOARCH)
		return
	}

	s := store.NewMemoryStore()

	fmt.Println("Start HTTP-Server listening on " + options.HttpBind)
	go RunHttpServer(options.HttpBind, s)

	fmt.Println("Start SMTP-server listening on " + options.SmtpBind)
	go RunSmtpServer(options.SmtpBind, s)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case s := <-c:
			if s == os.Interrupt {
				ticker.Stop()
				fmt.Println("\nExit. Bye!")
				return
			}
		case t := <-ticker.C:
			s.PurgeBefore(t.Add(-1 * options.Retention))
		}

	}
}
