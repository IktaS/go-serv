package serv

import (
	"html/template"
	"io"

	"github.com/IktaS/go-serv/pkg/serv"
)

// // GenerateServiceClient will create a service for client
// func GenerateServiceClient(s *Service) ([]byte, error) {

// }

// // GenerateMessageClient will create a Message for client
// func GenerateMessageClient(m *Message) ([]byte, error) {

// }

// // GenerateServClient will create a .ino for client
// func GenerateServClient(s *Serv) ([]byte, error) {

// }

// GenerateServServer will create a .go for server
func GenerateServServer(w io.Writer, s *serv.Gserv) error {
	tmpl, err := template.ParseFiles("./templates/server/base.tmpl", "./templates/server/message.tmpl", "./templates/server/service.tmpl")
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, s)
	if err != nil {
		return err
	}
	return nil

}
