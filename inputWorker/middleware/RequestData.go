package middleware

// Request
type CitizenAuthData struct {
	Otp string `json:"otp"`
}
type CustomFields struct {
	ObjectCode  string `json:"objectCode"`
	OrgCode     string `json:"orgCode"`
	OrgName     string `json:"orgName"`
	OrgPassword string `json:"orgPassword"`
	OrgToken    string `json:"orgToken"`
}

type AuthData struct{}
type SignData struct{}

type RequestBody struct {
	ServiceCode     string          `json:"serviceCode"`
	CitizenAuthData CitizenAuthData `json:"citizenAuthData"`
	CustomFields    CustomFields    `json:"customFields"`
	AuthData        AuthData        `json:"authData"`
	SignData        SignData        `json:"signData"`
}
