package email

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"

	"github.com/joho/godotenv"
)

type (
	email struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Html    string `json:"text"`
	}

	templateData struct {
		Domain    string
		FirstName string
		Link      string
	}
)

var uri = ""

func InitServer() {
	godotenv.Load(".env")
	uri = "https://api.mailgun.net/v3/" + os.Getenv("MAILGUN_DOMAIN") + "/messages"
}

func send(req *http.Request) error {
	client := &http.Client{}
	req.SetBasicAuth("api", os.Getenv("MAILGUN_API"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	bodybytes, _ := ioutil.ReadAll(res.Body)
	if err != nil || res.StatusCode > 299 {
		return fmt.Errorf("Error sending email: %s %s", err, bodybytes)
	}
	return nil
}

func SendBotAddLink(to, link, firsName string) error {
	tmpl := template.Must(template.ParseFiles("./internal/email/adminTemplate.html"))
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, templateData{os.Getenv("SERVER_DOMAIN"), firsName, link})

	v := url.Values{}
	v.Set("from", os.Getenv("MAILGUN_FROM"))
	v.Set("to", os.Getenv("ADMIN_EMAIL"))
	v.Set("subject", "Link to add Bot to server")
	v.Set("html", buf.String())

	req, _ := http.NewRequest("POST", uri, strings.NewReader(v.Encode()))
	return send(req)

}

func SendSheetsConsent(to, link, firsName string) error {
	tmpl := template.Must(template.ParseFiles("./internal/email/adminTemplate.html"))
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, templateData{os.Getenv("SERVER_DOMAIN"), firsName, link})

	v := url.Values{}
	v.Set("from", os.Getenv("MAILGUN_FROM"))
	v.Set("to", os.Getenv("ADMIN_EMAIL"))
	v.Set("subject", "Authorize Bot to access your sheets")
	v.Set("html", buf.String())

	req, _ := http.NewRequest("POST", uri, strings.NewReader(v.Encode()))
	return send(req)

}

func SendInviteLink(to, link, firsName string) error {
	tmpl := template.Must(template.ParseFiles("./internal/email/template.html"))
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, templateData{os.Getenv("SERVER_DOMAIN"), firsName, link})

	v := url.Values{}
	v.Set("from", os.Getenv("MAILGUN_FROM"))
	v.Set("to", to)
	v.Set("subject", `📣Your "Mind Unbound Trading" Discord Invitation!`)
	v.Set("html", buf.String())
	//pass the values to the request's body
	req, _ := http.NewRequest("POST", uri, strings.NewReader(v.Encode()))
	return send(req)

}
