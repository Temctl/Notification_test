package model

type PushNotificationModel struct {
	Title    string
	Body     string
	ImageUrl string
	Data     map[string]string
	Regnum   string
	CivilId  string
}

type XypNotification struct {
	Date        string
	ServiceName string
	ServiceDesc string
	OrgName     string
	Regnum      string
	CivilId     string
	RequestId   string
	ResultCode  int
	ClientId    int
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
	Type       NotificationType
	Regnum     string
	CivilId    string
	ExpireDate string
	Content    string
}

type RegularNotification struct {
	Content string
	Regnum  string
	CivilId string
}

type GroupNotification struct {
	IsAll    bool
	Regnums  []string
	CivilIds []string
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
