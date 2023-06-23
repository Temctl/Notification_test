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
	ContentData    XypContent `json:"contentData"`
	Regnum         string     `json:"regnum"`
	OperatorRegnum string     `json:"operatorRegnum"`
	CivilId        string     `json:"civilId"`
	ClientId       int        `json:"clientId"`
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

type Notifs struct {
	Id                int    `json:"id"`
	Org_id            string `json:"org_id"`
	Ws_operation_name string `json:"ws_operation_name"`
	Description       string `json:"description"`
	Created_at        string `json:"created_at"`
	Updated_at        string `json:"updated_at"`
	Reason            string `json:"reason"`
}

type OrginizationJson struct {
	Org_id        string   `json:"org_id"`
	Org_phone     string   `json:"org_phone"`
	Contact_phone string   `json:"contact_phone"`
	Contact_email string   `json:"contact_email"`
	Contact_web   string   `json:"contact_web"`
	Contract_id   string   `json:"contract_id"`
	Org_type_id   string   `json:"org_type_id"`
	Notifs        []Notifs `json:"notifs"`
}

type OrgInfoModel struct {
	Success        bool             `json:"success"`
	Client_is_soap string           `json:"client_is_soap"`
	Organization   OrginizationJson `json:"organization"`
}
