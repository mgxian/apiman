package mail

import (
	"crypto/tls"
	"strconv"

	//log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/pkg/setting"
	"gopkg.in/gomail.v2"
)

var (
	host, username, password, from string
	port                           int
	useTls                         bool
	mailer                         *gomail.Dialer
)

func loadConfig() {
	sec := setting.Cfg.Section("mail")
	host = sec.Key("host").String()
	port, _ = strconv.Atoi(sec.Key("port").String())
	username = sec.Key("username").String()
	password = sec.Key("password").String()
	from = sec.Key("from").String()
	useTls = sec.Key("tls").MustBool(false)
}

func MailInit() {
	loadConfig()
	mailer = gomail.NewDialer(host, port, username, password)
	mailer.TLSConfig = &tls.Config{InsecureSkipVerify: useTls}

}

func SendText(to []string, subject, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "apiman<"+username+">")
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text", content)

	return mailer.DialAndSend(m)
}

func SendHtml(to []string, subject, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "apiman<"+username+">")
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	return mailer.DialAndSend(m)
}
