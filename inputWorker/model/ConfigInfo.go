package model

type ConfigInfo struct {
	CivilId int    `json:"civilId"`
	Email   string `json:"email"`
}

type ApiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
