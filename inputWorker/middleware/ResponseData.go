package middleware

// Response
type Listdata struct {
	CivilId  string `json:"civilId"`
	No       string `json:"no"`
	Passport string `json:"passport"`
}

type ResponseData struct {
	DateOfExpiry string     `json:"dateOfExpiry"`
	Listdata     []Listdata `json:"listdata"`
}

type ResponseBody struct {
	Result          bool         `json:"result"`
	Message         string       `json:"message"`
	ResultCode      int          `json:"resultCode"`
	RequestId       string       `json:"requestId"`
	RequestData     struct{}     `json:"requestData"`
	Data            ResponseData `json:"data"`
	CitizenData     string       `json:"citizenData"`
	LegalEntityData string       `json:"legalEntityData"`
	ServiceCode     string       `json:"serviceCode"`
	CertificateDate string       `json:"certificateDate"`
}
