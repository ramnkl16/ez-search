package ezsmtp

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"text/template"

	"github.com/ramnkl16/ez-search/coredb"
	"github.com/ramnkl16/ez-search/logger"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}
func SendResetEmail(to, subject string, data interface{}) error {
	tmpCont, err := coredb.GetValue(coredb.Defaultbucket, bodyKey)
	if err != nil {
		logger.Error("Failed while get key", err)
		return err
	}
	subStr, err := coredb.GetValue(coredb.Defaultbucket, subjectKey)

	t, err := template.New("email").Parse(string(tmpCont))
	if err != nil {
		logger.Error("Failed while parse template", err)
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", string(subStr), mimeHeaders)))

	t.Execute(&body, data)

	// Sending email.
	err = smtp.SendMail(Conf.SmtpHost+":"+Conf.SmtpPort, getAuth(), Conf.From, []string{to}, body.Bytes())
	if err != nil {
		logger.Error("Failed while send email", err)
		return err
	}
	logger.Info("Email Sent!")
	return nil

}
func getAuth() smtp.Auth {
	conn, err := net.Dial("tcp", "smtp.office365.com:587")
	if err != nil {
		println(err)
	}

	c, err := smtp.NewClient(conn, Conf.SmtpHost)
	if err != nil {
		println(err)
	}

	tlsconfig := &tls.Config{
		ServerName: Conf.SmtpHost,
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		println(err)
	}

	auth := LoginAuth(Conf.From, Conf.Password)

	if err = c.Auth(auth); err != nil {
		println(err)
	}
	return auth

}

func SendEmailUsinggoogle(to []string) {

	// smtp server configuration.
	//smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}

	// Message.
	message := []byte("This is a really unimaginative message, I know.")

	// Authentication.
	auth := smtp.PlainAuth("", Conf.From, Conf.Password, Conf.SmtpHost)

	// Sending email.
	err := smtp.SendMail(fmt.Sprintf("%s:%s", Conf.SmtpHost, Conf.SmtpPort), auth, Conf.From, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
}
