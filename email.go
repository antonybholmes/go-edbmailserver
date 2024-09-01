package main

import (
	"bytes"
	"net/mail"
	"net/url"
	"strings"
	"text/template"

	"github.com/antonybholmes/go-edb-server-mailer/consts"
	"github.com/antonybholmes/go-mailer"

	"github.com/antonybholmes/go-mailer/mailserver"
)

const JWT_PARAM = "jwt"
const URL_PARAM = "url"

type EmailBody struct {
	Name       string
	From       string
	Time       string
	Link       string
	DoNotReply string
}

func SendPasswordlessSigninEmail(qe *mailer.RedisQueueEmail) error {

	var file string

	if qe.CallBackUrl != "" {
		file = "templates/email/passwordless/web.html"
	} else {
		file = "templates/email/passwordless/api.html"
	}

	go SendEmailWithToken("Passwordless Sign In",
		qe,
		file)

	return nil
}

func SendEmailWithToken(subject string,
	qe *mailer.RedisQueueEmail,
	file string) error {

	address, err := mail.ParseAddress(qe.To)

	if err != nil {
		return err
	}

	var body bytes.Buffer

	t, err := template.ParseFiles(file)

	if err != nil {
		return err
	}

	var firstName string = ""

	if qe.Name != "" {
		firstName = qe.Name
	} else {
		firstName = strings.Split(address.Address, "@")[0]
	}

	firstName = strings.Split(firstName, " ")[0]

	if qe.CallBackUrl != "" {
		callbackUrl, err := url.Parse(qe.CallBackUrl)

		if err != nil {
			return err
		}

		params, err := url.ParseQuery(callbackUrl.RawQuery)

		if err != nil {
			return err
		}

		if qe.VisitUrl != "" {
			params.Set(URL_PARAM, qe.VisitUrl)
		}

		params.Set(JWT_PARAM, qe.Token)

		callbackUrl.RawQuery = params.Encode()

		link := callbackUrl.String()

		err = t.Execute(&body, EmailBody{
			Name:       firstName,
			Link:       link,
			From:       consts.NAME,
			Time:       qe.Ttl,
			DoNotReply: consts.DO_NOT_REPLY,
		})

		if err != nil {
			return err
		}
	} else {
		err = t.Execute(&body, EmailBody{
			Name:       firstName,
			Link:       qe.Token,
			From:       consts.NAME,
			Time:       qe.Ttl,
			DoNotReply: consts.DO_NOT_REPLY,
		})

		if err != nil {
			return err
		}
	}

	err = mailserver.SendHtmlEmail(address, subject, body.String())

	if err != nil {
		return err
	}

	return nil
}
