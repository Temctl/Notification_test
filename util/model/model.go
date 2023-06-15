package model

type XypNotification struct {
	Date        string `json:"date"`
	ServiceName string `json:"serviceName"`
	ServiceDesc string `json:"serviceDesc"`
	OrgName     string `json:"orgName"`
	Regnum      string `json:"regnum"`
	CivilId     string `json:"civilId"`
	RequestId   string `json:"requestId"`
	ResultCode  int    `json:"resultCode"`
	ClientId    int    `json:"clientId"`
}

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
	Content    string           `json:"content"`
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

type RegularNotification struct {
	Content string `json:"content"`
	Regnum  string `json:"regnum"`
	CivilId string `json:"civilId"`
}

type GroupNotification struct {
	IsAll    bool     `json:"isAll"`
	Content  string   `json:"content"`
	Regnums  []string `json:"regnums"`
	CivilIds []string `json:"civilIds"`
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
