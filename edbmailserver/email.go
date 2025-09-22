package edbmailserver

import (
	"bytes"
	"html/template"
	"net/mail"
	"net/url"
	"strings"

	"github.com/antonybholmes/go-edbmailserver/consts"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/rs/zerolog/log"

	"github.com/antonybholmes/go-mailserver/sesmailserver"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const JWT_PARAM = "jwt"
const FOOTER_TEMPLATE = "templates/email/footer.html"

//const REDIRECT_URL_PARAM = "redirectUrl"

type TimeSensitive struct {
	ContentType string
	Time        string
}
type EmailBody struct {
	Name          string
	From          string
	TimeSensitive *TimeSensitive
	Payload       *mailserver.Payload
	Link          string
	DoNotReply    string
}

func SendPasswordlessSigninEmail(mail *mailserver.MailItem) error {

	var file string

	if mail.Mode == "api" {
		file = "templates/email/passwordless/api.html"
	} else {
		file = "templates/email/passwordless/web.html"
	}

	//file := "templates/email/passwordless/web.html"

	return SendEmailWithToken("Passwordless Sign In",
		mail,
		consts.URL_SIGN_IN,
		file)
}

func SendVerifyEmail(mail *mailserver.MailItem) error {

	var file string

	if mail.Mode == "api" {
		file = "templates/email/verify/api.html"
	} else {
		file = "templates/email/verify/web.html"
	}

	return SendEmailWithToken("Email Address Verification",
		mail,
		consts.URL_VERIFY_EMAIL,
		file)
}

func SendVerifiedEmail(mail *mailserver.MailItem) error {

	file := "templates/email/verify/verified.html"

	return SendEmailWithToken("Email Address Verified",
		mail,
		"",
		file)
}

func SendPasswordResetEmail(mail *mailserver.MailItem) error {

	var file string

	if mail.LinkUrl != "" {
		file = "templates/email/password/reset/web.html"
	} else {
		file = "templates/email/password/reset/api.html"
	}

	return SendEmailWithToken("Password Reset",
		mail,
		mail.LinkUrl,
		file)
}

func SendPasswordUpdatedEmail(mail *mailserver.MailItem) error {

	var file string

	if mail.LinkUrl != "" {
		file = "templates/email/password/switch-to-passwordless.html"
	} else {
		file = "templates/email/password/updated.html"
	}

	return SendEmailWithToken("Password Updated",
		mail,
		mail.LinkUrl,
		file)
}

func SendEmailResetEmail(mail *mailserver.MailItem) error {

	var file string

	if mail.LinkUrl != "" {
		file = "templates/email/email/reset/web.html"
	} else {
		file = "templates/email/email/reset/api.html"
	}

	return SendEmailWithToken("Email Reset",
		mail,
		mail.LinkUrl,
		file)
}

func SendEmailUpdatedEmail(mail *mailserver.MailItem) error {

	file := "templates/email/email/updated.html"

	return SendEmailWithToken("Email Updated",
		mail,
		mail.LinkUrl,
		file)
}

func SendAccountCreatedEmail(mail *mailserver.MailItem) error {

	file := "templates/email/account/created.html"

	return SendEmailWithToken("Account Created",
		mail,
		mail.LinkUrl,
		file)
}

func SendAccountUpdatedEmail(mail *mailserver.MailItem) error {

	file := "templates/email/account/updated.html"

	return SendEmailWithToken("Account Updated",
		mail,
		mail.LinkUrl,
		file)
}

func SendOTPEmail(mail *mailserver.MailItem) error {

	file := "templates/email/otp/otp.html"

	//log.Debug().Msgf("send totp email to %s", mail.To)

	err := SendEmailWithToken("One-Time Passcode For Experiments Sign In",
		mail,
		"",
		file)

	// if err != nil {
	// 	log.Debug().Msgf("error sending totp email to %s: %v", mail.To, err)
	// }

	return err
}

func SendEmailWithToken(subject string,
	m *mailserver.MailItem,
	linkUrl string,
	file string) error {

	//baseName := filepath.Base(file)

	//log.Debug().Msgf("send email %s to %s using %s", subject, m.To, baseName)

	address, err := mail.ParseAddress(m.To)

	if err != nil {
		return err
	}

	var body bytes.Buffer

	t, err := template.ParseFiles(file, FOOTER_TEMPLATE)

	if err != nil {
		return err
	}

	var firstName string = ""

	if m.Name != "" {
		firstName = m.Name
	} else {
		firstName = strings.Split(address.Address, "@")[0]
		c := cases.Title(language.English)
		firstName = c.String(firstName)
	}

	firstName = strings.Split(firstName, " ")[0]

	link := ""

	if m.Payload != nil && m.Payload.DataType == "link" {
		link = m.Payload.Data
	}

	if linkUrl != "" {
		parsedUrl, err := url.Parse(linkUrl)

		if err != nil {
			return err
		}

		params, err := url.ParseQuery(parsedUrl.RawQuery)

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

		if m.Payload != nil && m.Payload.DataType == "jwt" {
			params.Set(JWT_PARAM, m.Payload.Data)
		}

		// once we've added extra params, update the
		// raw query again
		parsedUrl.RawQuery = params.Encode()

		// the complete url with params
		link = parsedUrl.String()

		// err = t.Execute(&body, EmailBody{
		// 	Name:       firstName,
		// 	Link:       link,
		// 	From:       consts.NAME,
		// 	Time:       m.TTL,
		// 	DoNotReply: consts.DO_NOT_REPLY,
		// })

		// if err != nil {
		// 	return err
		// }
	}

	emailBody := EmailBody{
		Name:       firstName,
		Payload:    m.Payload,
		Link:       link,
		From:       consts.NAME,
		DoNotReply: consts.DO_NOT_REPLY,
	}

	if m.TTL != "" {
		contentType := "token"

		if link != "" {
			contentType = "link"
		} else {
			if m.Payload != nil {
				if m.Payload.DataType != "jwt" {
					contentType = m.Payload.DataType
				}
			}
		}

		emailBody.TimeSensitive = &TimeSensitive{
			ContentType: contentType,
			Time:        m.TTL,
		}
	}

	err = t.ExecuteTemplate(&body, "layout", emailBody)

	if err != nil {
		return err
	}

	//log.Debug().Msgf("email %v %v %v", address, subject, body.String())

	err = sesmailserver.SendHtmlMail(address, subject, body.String())

	if err != nil {
		return err
	}

	log.Info().Msgf("Email of type %s sent to %s for %s", m.EmailType, address.Address, subject)

	return nil
}
