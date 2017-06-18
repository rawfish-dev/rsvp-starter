package security

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
)

const (
	googleReCAPTCHAVerifyURL = "https://www.google.com/recaptcha/api/siteverify"
	googleReCAPTCHASecret    = "6Ld9PQkUAAAAAE6oyVaT8B1XTcOYlpeOOYlyiUTX"
)

type googleReCAPTCHAVerifyResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

func (s *service) VerifyReCAPTCHA(token string) (valid bool) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	fullURL, err := url.Parse(googleReCAPTCHAVerifyURL)
	if err != nil {
		ctxLogger.Errorf("security service - unable to parse reCAPTCHA verify url due to %v", err)
		return false
	}
	urlQuery := fullURL.Query()
	urlQuery.Add("secret", googleReCAPTCHASecret)
	urlQuery.Add("response", token)
	fullURL.RawQuery = urlQuery.Encode()

	verifyReq, err := http.NewRequest("POST", fullURL.String(), nil)
	if err != nil {
		ctxLogger.Errorf("security service - unable to create reCAPTCHA verify request due to %v", err)
		return false
	}

	client := &http.Client{}
	resp, err := client.Do(verifyReq)
	if err != nil {
		ctxLogger.Errorf("security service - unable to complete reCAPTCHA verify due to %v", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctxLogger.Errorf("security service - unable to read reCAPTCHA response body due to %v", err)
		return false
	}

	var verifyResp googleReCAPTCHAVerifyResponse

	err = json.Unmarshal(body, &verifyResp)
	if err != nil {
		ctxLogger.Errorf("security service - unable to unwrap reCAPTCHA response body due to %v", err)
		return false
	}

	if len(verifyResp.ErrorCodes) > 0 {
		ctxLogger.Errorf("security service - validation of reCAPTCHA failed due to %v", verifyResp.ErrorCodes)
		return false
	}

	return verifyResp.Success
}
