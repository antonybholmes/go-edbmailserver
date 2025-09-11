package main

import (
	"bytes"
	"net/mail"
	"net/url"
	"strings"
	"text/template"

	"github.com/antonybholmes/go-edbmailserver/consts"
	mailserver "github.com/antonybholmes/go-mailserver"

	"github.com/antonybholmes/go-mailserver/sesmailserver"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const JWT_PARAM = "jwt"

//const REDIRECT_URL_PARAM = "redirectUrl"

type EmailBody struct {
	Name       string
	From       string
	Time       string
	Link       string
	DoNotReply string
}

func SendPasswordlessSigninEmail(qe *mailserver.QueueEmail) error {

	var file string

	if qe.Mode == "api" {
		file = "templates/email/passwordless/api.html"
	} else {
		file = "templates/email/passwordless/web.html"
	}

	//file := "templates/email/passwordless/web.html"

	return SendEmailWithToken("Passwordless Sign In",
		qe,
		consts.URL_SIGN_IN,
		file)
}

func SendVerifyEmail(qe *mailserver.QueueEmail) error {

	var file string

	if qe.Mode == "api" {
		file = "templates/email/verify/api.html"
	} else {
		file = "templates/email/verify/web.html"
	}

	return SendEmailWithToken("Email Address Verification",
		qe,
		consts.URL_VERIFY_EMAIL,
		file)
}

func SendVerifiedEmail(qe *mailserver.QueueEmail) error {

	file := "templates/email/verify/verified.html"

	return SendEmailWithToken("Email Address Verified",
		qe,
		"",
		file)
}

func SendPasswordResetEmail(qe *mailserver.QueueEmail) error {

	var file string

	if qe.LinkUrl != "" {
		file = "templates/email/password/reset/web.html"
	} else {
		file = "templates/email/password/reset/api.html"
	}

	return SendEmailWithToken("Password Reset",
		qe,
		qe.LinkUrl,
		file)
}

func SendPasswordUpdatedEmail(qe *mailserver.QueueEmail) error {

	var file string

	if qe.LinkUrl != "" {
		file = "templates/email/password/switch-to-passwordless.html"
	} else {
		file = "templates/email/password/updated.html"
	}

	return SendEmailWithToken("Password Updated",
		qe,
		qe.LinkUrl,
		file)
}

func SendEmailResetEmail(qe *mailserver.QueueEmail) error {

	var file string

	if qe.LinkUrl != "" {
		file = "templates/email/email/reset/web.html"
	} else {
		file = "templates/email/email/reset/api.html"
	}

	return SendEmailWithToken("Email Reset",
		qe,
		qe.LinkUrl,
		file)
}

func SendEmailUpdatedEmail(qe *mailserver.QueueEmail) error {

	file := "templates/email/email/updated.html"

	return SendEmailWithToken("Email Updated",
		qe,
		qe.LinkUrl,
		file)
}

func SendAccountCreatedEmail(qe *mailserver.QueueEmail) error {

	file := "templates/email/account/created.html"

	return SendEmailWithToken("Account Created",
		qe,
		qe.LinkUrl,
		file)
}

func SendAccountUpdatedEmail(qe *mailserver.QueueEmail) error {

	file := "templates/email/account/updated.html"

	return SendEmailWithToken("Account Updated",
		qe,
		qe.LinkUrl,
		file)
}

func SendOTPEmail(qe *mailserver.QueueEmail) error {

	file := "templates/email/otp/totp.html"

	//log.Debug().Msgf("send totp email to %s", qe.To)

	err := SendEmailWithToken("One-Time Password (OTP)",
		qe,
		"",
		file)

	// if err != nil {
	// 	log.Debug().Msgf("error sending totp email to %s: %v", qe.To, err)
	// }

	return err
}

func SendEmailWithToken(subject string,
	qe *mailserver.QueueEmail,
	linkUrl string,
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
		c := cases.Title(language.English)
		firstName = c.String(firstName)
	}

	firstName = strings.Split(firstName, " ")[0]

	if linkUrl != "" {
		_linkUrl, err := url.Parse(linkUrl)

		if err != nil {
			return err
		}

		params, err := url.ParseQuery(_linkUrl.RawQuery)

		if err != nil {
			return err
		}

		// after the callback url does some validation, we can
		// goto a different url to make the user experience
		// better
		// this feature is mostly unused since the visit url
		// is normally encoded in the attached jwt to prevent
		// tampering
		//if qe.RedirectUrl != "" {
		//	params.Set(REDIRECT_URL_PARAM, qe.RedirectUrl)
		//}

		if qe.Token != "" {
			params.Set(JWT_PARAM, qe.Token)
		}

		// once we've added extra params, update the
		// raw query again
		_linkUrl.RawQuery = params.Encode()

		// the complete url with params
		link := _linkUrl.String()

		err = t.Execute(&body, EmailBody{
			Name:       firstName,
			Link:       link,
			From:       consts.NAME,
			Time:       qe.TTL,
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
			Time:       qe.TTL,
			DoNotReply: consts.DO_NOT_REPLY,
		})

		if err != nil {
			return err
		}
	}

	//log.Debug().Msgf("email %v %v %v", address, subject, body.String())

	err = sesmailserver.SendHtmlEmail(address, subject, body.String())

	if err != nil {
		return err
	}

	return nil
}
