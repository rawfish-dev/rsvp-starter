package twilio

type TwilioServiceProvider interface {
	SendSMS(mobilePhoneNumber, message string) (success bool, err error)
}
