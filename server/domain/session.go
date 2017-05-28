package domain

type SessionCreateRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	ReCAPTCHAToken string `json:"reCAPTCHA"`
}

type SessionCreateResponse struct {
	Username  string `json:"username"`
	AuthToken string `json:"authToken"`
}
