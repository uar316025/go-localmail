package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/emersion/go-imap/server"
	"github.com/emersion/go-smtp"
	"github.com/iii-xvi/go-localmail"
	"log"
	"os"
	"os/signal"
	"time"
)

type options struct {
	ImapPort        int
	SmtpPort        int
	HttpPort        int
	CertificatePath string
}

func parseCliArgs() *options {
	args := new(options)

	flag.IntVar(&args.ImapPort, "imap", 2993, "port number the IMAP server will listen on")
	flag.IntVar(&args.SmtpPort, "smtp", 2026, "port number the SMTP server will listen on")
	//flag.IntVar(&args.HttpPort, "http", 2080, "port number the HTTP server will listen on")
	flag.StringVar(&args.CertificatePath, "cert", "", "path to SSL certificate .crt (."+
		"key should be in same dir with same name)")

	flag.Parse()
	return args
}

func main() {
	var cert tls.Certificate

	args := parseCliArgs()
	//servicesDone := &sync.WaitGroup{}

	// load cert if needed
	tlsEnabled := args.CertificatePath != ""
	if tlsEnabled {
		var err error
		keyFile := args.CertificatePath[:len(args.CertificatePath)-4] + ".key"
		cert, err = tls.LoadX509KeyPair(args.CertificatePath, keyFile)
		if err != nil {
			log.Println(err)
			return
		}
	}

	// init IMAP
	// Create a smtpMemory backend
	backend := localmail.NewBackend()

	// Create a new server
	sImap := server.New(backend.IMAP())
	sImap.Addr = fmt.Sprintf("0.0.0.0:%d", args.ImapPort)
	// Since we will use this server for testing only, we can allow plain text
	// authentication over unencrypted connections
	sImap.AllowInsecureAuth = true

	// init SMTP
	sSmtp := smtp.NewServer(backend.SMTP())

	sSmtp.Addr = fmt.Sprintf(":%d", args.SmtpPort)
	sSmtp.Domain = "0.0.0.0"
	sSmtp.ReadTimeout = 10 * time.Second
	sSmtp.WriteTimeout = 10 * time.Second
	sSmtp.MaxMessageBytes = 1024 * 1024 * 24 // 24 MB
	sSmtp.MaxRecipients = 50
	sSmtp.AllowInsecureAuth = true

	go func() {
		var err error
		if tlsEnabled {
			log.Println("Starting SMTP+TLS server at " + sSmtp.Domain + sSmtp.Addr)
			sSmtp.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
			err = sSmtp.ListenAndServeTLS()
		} else {
			log.Println("Starting SMTP server at " + sSmtp.Domain + sSmtp.Addr)
			err = sSmtp.ListenAndServe()
		}
		if err != nil {
			log.Println(err)
		}
	}()
	go func() {
		var err error
		if tlsEnabled {
			log.Println("Starting IMAP+TLS server at " + sImap.Addr)
			sSmtp.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
			err = sImap.ListenAndServeTLS()
		} else {
			log.Println("Starting IMAP server at " + sImap.Addr)
			err = sImap.ListenAndServe()
		}
		if err != nil {
			log.Println(err)
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Println("Press Ctrl+C to stop.")
	// Waiting for SIGINT (pkill -2)
	<-stop

	// stop all
	if err := sImap.Close(); err != nil {
		log.Println(err)
	}
	sSmtp.Close()
}
