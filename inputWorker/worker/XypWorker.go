package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/Temctl/E-Notification/util/model"
)

type XypNotifMarshal struct {
	Regnum         string `json:"regnum"`
	OperatorRegnum string `json:"operatorRegnum"`
	Date           string `json:"date"`
	ServiceName    string `json:"serviceName"`
	ServiceDesc    string `json:"serviceDesc"`
	OrgName        string `json:"orgName"`
	RequestId      string `json:"requestId"`
	ResultCode     int    `json:"resultCode"`
	ResultMessage  string `json:"resultMessage"`
	ClientId       int    `json:"clientId"`
	CivilId        string `json:"civilId"`
}

func XypWorker() {
	elog.Info().Println("XYP NOTIF WORKER STARTED...")
	// ----------------------------------------------------------------------
	// REDIS ----------------------------------------------------------------
	// ----------------------------------------------------------------------

	redisClient, err := connections.ConnectionRedis()
	if err != nil {
		elog.Error().Panic(err)
	}

	// ----------------------------------------------------------------------
	// DB CONNECTION --------------------------------------------------------
	// ----------------------------------------------------------------------
	collection, client, err := connections.GetMongoCollection(model.XYPNOTIFICATION)
	if err != nil {
		elog.Error().Panic(err)
	}
	defer client.Disconnect(context.Background())
	// ------------------------------------------------------------
	// Infinite loop to continuously pop items from the list ------
	// ------------------------------------------------------------
	for {
		// --------------------------------------------------------
		// Pop an item from the list using the BLPOP command ------
		// --------------------------------------------------------
		result, err := redisClient.BLPop(0, "XYPNOTIF").Result()

		if err != nil {
			elog.Error().Println("Error:", err)
			continue
		}
		// --------------------------------------------------------
		// Check if an item was successfully popped ---------------
		// --------------------------------------------------------
		//
		// [queue:queue, value] len 2 bol zuv gesen ug
		// [queue:queue, value1]
		// [queue:queue, value2]

		xypNotifMarshal := XypNotifMarshal{}

		if len(result) == 2 {
			value := result[1]
			err := json.Unmarshal([]byte(value), &xypNotifMarshal)
			if err != nil {
				elog.Error().Println("Error:", err)
				continue
			}
			xypNotif := model.XypNotification{
				Regnum:         xypNotifMarshal.Regnum,
				OperatorRegnum: xypNotifMarshal.OperatorRegnum,
				CivilId:        xypNotifMarshal.CivilId,
				ClientId:       xypNotifMarshal.ClientId,
				ContentData: model.XypContent{
					OrgName:     xypNotifMarshal.OrgName,
					ServiceDesc: xypNotifMarshal.ServiceDesc,
					Date:        xypNotifMarshal.Date,
					ServiceName: xypNotifMarshal.ServiceName,
					RequestId:   xypNotifMarshal.RequestId,
					ResultCode:  xypNotifMarshal.ResultCode,
				},
			}
			// Insert the document
			_, insertErr := collection.InsertOne(context.Background(), xypNotif)
			if insertErr != nil {
				elog.Error().Panic(insertErr)
				continue
			}
			fmt.Println(xypNotif)
		}
	}
}

// ('queue:notification', '{"regnum":"уц03211011","operatorRegnum":"","date":"2023-06-15 16:39:52","serviceName":"WS100101_getCitizenIDCardInfo","serviceDesc":"Иргэний үнэмлэхний мэдээлэл дамжуулах сервис","orgName":"Сангийн яам","requestId":"a43b49c3-ba1e-4eb0-883e-bf04b2ecacfc","resultCode":0,"resultMessage":"амжилттай","clientId":126,"retryCount":0,"civilId":"891834062934"}')
