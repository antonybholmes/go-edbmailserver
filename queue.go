package main

const (
	QUEUE_EMAIL_CHANNEL = "email-channel"
)

const (
	QUEUE_EMAIL_TYPE_VERIFY           = "verify"
	QUEUE_EMAIL_TYPE_VERIFIED         = "verified"
	QUEUE_EMAIL_TYPE_PASSWORDLESS     = "passwordless"
	QUEUE_EMAIL_TYPE_PASSWORD_RESET   = "password-reset"
	QUEUE_EMAIL_TYPE_PASSWORD_UPDATED = "password-updated"
	QUEUE_EMAIL_TYPE_EMAIL_RESET      = "email-reset"
	QUEUE_EMAIL_TYPE_EMAIL_UPDATED    = "email-updated"
	QUEUE_EMAIL_TYPE_ACCOUNT_CREATED  = "account-created"
	QUEUE_EMAIL_TYPE_ACCOUNT_UPDATED  = "account-updated"
)

type QueueEmail struct {
	Name      string `json:"name"`
	To        string `json:"to"`
	EmailType string `json:"type"`
	Ttl       string `json:"ttl"`
	Token     string `json:"token"`
	//RedirectUrl string `json:"redirectUrl"`
	LinkUrl string `json:"linkUrl"`
	// Used to indicate if emails should be for web users or
	// geared towards API users who just want the tokens without
	// any links
	Mode string `json:"mode"`
}
