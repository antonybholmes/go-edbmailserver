package mail

const (
	EmailQueueTypeVerify          = "verify"
	EmailQueueTypeVerified        = "verified"
	EmailQueueTypePasswordless    = "passwordless"
	EmailQueueTypePasswordReset   = "password-reset"
	EmailQueueTypePasswordUpdated = "password-updated"
	EmailQueueTypeEmailReset      = "email-reset"
	EmailQueueTypeEmailUpdated    = "email-updated"
	EmailQueueTypeAccountCreated  = "account-created"
	EmailQueueTypeAccountUpdated  = "account-updated"
	EmailQueueTypeOTP             = "otp"
)
