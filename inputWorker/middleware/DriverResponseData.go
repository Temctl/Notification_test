package middleware

// Response
type DListdata struct {
	ExpirationDate string `json:"expirationDate"`
	Regnum         string `json:"regnum"`
}

type DResponseData struct {
	Listdata []DListdata `json:"listdata"`
}

type DResponseBody struct {
	Result          bool          `json:"result"`
	Message         string        `json:"message"`
	ResultCode      int           `json:"resultCode"`
	RequestId       string        `json:"requestId"`
	RequestData     struct{}      `json:"requestData"`
	Data            DResponseData `json:"data"`
	CitizenData     string        `json:"citizenData"`
	LegalEntityData string        `json:"legalEntityData"`
	ServiceCode     string        `json:"serviceCode"`
	CertificateDate string        `json:"certificateDate"`
}
