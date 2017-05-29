package twilio

import (
	"github.com/rawfish-dev/rsvp-starter/server/services/base"

	"github.com/sfreiberg/gotwilio"
)

const (
	accountSID = "AC4c71fb74e15ed9226cb99b48bb917efc"
	authToken  = "017f255d1ced3fc236361d5560dc5490"
	fromNumber = "+18778766276"
)

type service struct {
	baseService  *base.Service
	twilioClient *gotwilio.Twilio
}

func NewService(baseService *base.Service) TwilioServiceProvider {
	return &service{baseService, gotwilio.NewTwilioClient(accountSID, authToken)}
}

func (s *service) SendSMS(mobilePhoneNumber, message string) (success bool, err error) {
	smsResponse, exception, err := s.twilioClient.SendSMS(fromNumber, mobilePhoneNumber, message, "", "")
	if err != nil {
		s.baseService.Errorf("twilio service - unable to send SMS due to %v", err)
		return false, err
	}

	if smsResponse != nil {
		s.baseService.Infof("twilio service - SMS response %+v", smsResponse)
	}

	if exception != nil {
		s.baseService.Errorf("twilio service - unable to send SMS successfully due to %v", *exception)
		return false, err
	}

	return true, nil
}
