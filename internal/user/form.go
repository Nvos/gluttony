package user

type LoginForm struct {
	Credentials Credentials `json:"credentials"`
	RedirectURL string `json:"redirectUrl"`
}
