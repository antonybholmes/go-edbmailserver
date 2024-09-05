package main

import (
	"bytes"
	"net/mail"
	"net/url"
	"strings"
	"text/template"

	"github.com/antonybholmes/go-edb-server-mailer/consts"
	"github.com/antonybholmes/go-mailer"
	"github.com/rs/zerolog/log"

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

	log.Debug().Msgf("cheese")
	go SendEmailWithToken("Passwordless Sign In",
		qe,
		file)

	return nil
}

func SendVerifyEmail(qe *mailer.RedisQueueEmail) error {

	var file string

	if qe.CallBackUrl != "" {
		file = "templates/email/verify/web.html"
	} else {
		file = "templates/email/verify/api.html"
	}

	go SendEmailWithToken("Email Address Verification",
		qe,
		file)

	return nil
}

func SendVerifiedEmail(qe *mailer.RedisQueueEmail) error {

	file := "templates/email/verify/verified.html"

	go SendEmailWithToken("Email Address Verified",
		qe,
		file)

	return nil
}

func SendPasswordResetEmail(qe *mailer.RedisQueueEmail) error {

	var file string

	if qe.CallBackUrl != "" {
		file = "templates/email/password/reset/web.html"
	} else {
		file = "templates/email/password/reset/api.html"
	}

	go SendEmailWithToken("Password Reset",
		qe,
		file)

	return nil
}

func SendPasswordUpdatedEmail(qe *mailer.RedisQueueEmail) error {

	var file string

	if qe.CallBackUrl != "" {
		file = "templates/email/password/switch-to-passwordless.html"
	} else {
		file = "templates/email/password/updated.html"
	}

	go SendEmailWithToken("Password Updated",
		qe,
		file)

	return nil
}

func SendEmailResetEmail(qe *mailer.RedisQueueEmail) error {

	var file string

	if qe.CallBackUrl != "" {
		file = "templates/email/email/reset/web.html"
	} else {
		file = "templates/email/email/reset/api.html"
	}

	go SendEmailWithToken("Email Reset",
		qe,
		file)

	return nil
}

func SendEmailUpdatedEmail(qe *mailer.RedisQueueEmail) error {

	file := "templates/email/email/updated.html"

	go SendEmailWithToken("Email Updated",
		qe,
		file)

	return nil
}

func SendAccountCreatedEmail(qe *mailer.RedisQueueEmail) error {

	file := "templates/email/account/created.html"

	go SendEmailWithToken("Account Created",
		qe,
		file)

	return nil
}

func SendAccountUpdatedEmail(qe *mailer.RedisQueueEmail) error {

	file := "templates/email/account/updated.html"

	go SendEmailWithToken("Account Updated",
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

		if qe.Token != "" {
			params.Set(JWT_PARAM, qe.Token)
		}

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

	//log.Debug().Msgf("awhat %v", address)

	err = mailserver.SendHtmlEmail(address, subject, body.String())

	if err != nil {
		return err
	}

	return nil
}
