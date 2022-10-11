package main

import (
	"io"
	"io/ioutil"
	"log"
	"time"
	"flag"
	"strconv"
	"github.com/emersion/go-smtp"
)

type Backend struct{}

func (bkd *Backend) NewSession(_ smtp.ConnectionState, _ string) (smtp.Session, error) {
	return &Session{}, nil
}

func (bkd *Backend) AnonymousLogin(_ *smtp.ConnectionState) (smtp.Session, error) {
	return &Session{}, nil
}

func (bkd *Backend) Login(_ *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	log.Println("username:", username)
	log.Println("password:", password)
	return &Session{}, nil
}

type Session struct{}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	log.Println("Mail from:", from)
	return nil
}

func (s *Session) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
		log.Println("Data:", string(b))
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func main() {
	port := flag.Int("port", 2525, "The port on which it is being listened to.")
	address := flag.String("address", "127.0.0.1", "The address on which it is being listened to.")
	domain := flag.String("Domain", "ACME.local", "Email domain")
	flag.Parse()
	be := &Backend{}
	s := smtp.NewServer(be)
	s.Addr = *address + ":" + strconv.Itoa(*port)
	s.Domain = *domain
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true
	log.Println("SMTP credentials sniffer server. For more features run app with \"--help\" options")
	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
