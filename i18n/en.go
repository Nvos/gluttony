package i18n

var enBundle = Bundle{ //nolint:gochecknoglobals // initialized globally as readonly
	Lang: "en",
	Messages: map[string]string{
		"login.header":        "Sign in to Gluttony",
		"login.description":   "Enter your account credentials below",
		"login.form.username": "Username",
		"login.form.password": "Password",
		"login.form.submit":   "Login",
	},
}
