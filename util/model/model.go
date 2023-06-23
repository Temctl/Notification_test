package model

type XypContent struct {
	OrgName     string `json:"orgName"`
	ServiceDesc string `json:"serviceDesc"`
	Date        string `json:"date"`
	ServiceName string `json:"serviceName"`
	RequestId   string `json:"requestId"`
	ResultCode  int    `json:"resultCode"`
}

type XypNotification struct {
	ContentData    []XypContent `json:"contentData"`
	Regnum         string       `json:"regnum"`
	OperatorRegnum string       `json:"operatorRegnum"`
	CivilId        string       `json:"civilId"`
	ClientId       int          `json:"clientId"`
}

type Collections string

const (
	XYPNOTIFICATION       = "xypnotification"
	ATTENTIONNOTIFICATION = "attentionnotification"
	OUTLOG                = "outlog"
	INLOG                 = "inlog"
)

type NotificationType string

const (
	DEFAULT = iota
	DRIVERLICENSEEXPIRE30
	IDCARDGOINGTOEXPIRE
	DRIVERLICENSEEXPIRED
	INTPASSPORTGOINGTOEXPIRE
)

type AttentionNotification struct {
	Type       NotificationType `json:"type"`
	Regnum     string           `json:"regnum"`
	CivilId    string           `json:"civilId"`
	ExpireDate string           `json:"expireDate"`
	Passport   string           `json:"passport"`
}

type PushNotificationModel struct {
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	ImageUrl string            `json:"imageUrl"`
	Data     map[string]string `json:"data"`
	Regnum   string            `json:"regnum"`
	CivilId  string            `json:"civilId"`
	Type     NotificationType  `json:"type"`
}

type EmailModel struct {
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	Destination string `json:"destination"`
	From        string `json:"from"`
	Regnum      string `json:"regnum"`
	CivilId     string `json:"civilId"`
}

type MessengerModel struct {
	Body    string `json:"body"`
	Regnum  string `json:"regnum"`
	CivilId string `json:"civilId"`
}

type RegularNotificationModel struct {
	IsAll    bool              `json:"isAll"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	ImageUrl string            `json:"imageUrl"`
	Data     map[string]string `json:"data"`
	Regnums  []string          `json:"regnums"`
	CivilIds []string          `json:"civilIds"`
	Tokens   []string          `json:"tokens"`
	Type     NotificationType  `json:"type"`
}

type UserConfigNotification struct {
	CivilId         string `json:"civilId"`
	Regnum          string `json:"regnum"`
	EmailAddress    string `json:"emailAddress"`
	IsSms           bool   `json:"isSms"`
	IsEmail         bool   `json:"isEmail"`
	IsPush          bool   `json:"isPush"`
	IsNationalEmail bool   `json:"isNationalEmail"`
	Social          bool   `json:"social"`
}
