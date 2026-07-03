package output

type LoginOutput struct {
	Tokens `json:"tokens"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewLoginOutput(accessToken, refreshToken string) *LoginOutput {
	return &LoginOutput{Tokens{AccessToken: accessToken, RefreshToken: refreshToken}}
}
