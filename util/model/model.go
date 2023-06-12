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

type AttentionType string

const (
	TYPENOTSPECIFIED = iota
	DRIVERLICENSEEXPIRE30
	IDCARDGOINGTOEXPIRE
	DRIVERLICENSEEXPIRED
	INTPASSPORTGOINGTOEXPIRE
)

type AttentionNotification struct {
	Type       AttentionType
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

type RedisConfigModel struct {
	CivilId         string
	Regnum          string
	EmailAddress    string
	IsSms           bool
	IsEmail         bool
	IsPush          bool
	IsNationalEmail bool
	Social          bool
}
