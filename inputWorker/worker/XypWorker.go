package worker

import (
	"encoding/json"
	"fmt"

	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/Temctl/E-Notification/util/model"
)

func XypWorker() {
	elog.Info().Println("XYP NOTIF WORKER STARTED...")
	// ----------------------------------------------------------------------
	// REDIS ----------------------------------------------------------------
	// ----------------------------------------------------------------------

	redisClient, err := connections.ConnectionRedis()
	if err != nil {
		elog.Error().Panic(err)
	}

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

		xypNotif := model.XypNotification{}

		if len(result) == 2 {
			value := result[1]
			err := json.Unmarshal([]byte(value), &xypNotif)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(xypNotif.CivilId)
			fmt.Println(xypNotif.ClientId)
			fmt.Println(xypNotif.Date)
			fmt.Println(xypNotif.OrgName)
			fmt.Println(xypNotif.Regnum)
			fmt.Println(xypNotif.RequestId)
			fmt.Println(xypNotif.ResultCode)
			fmt.Println(xypNotif.ServiceDesc)
			fmt.Println(xypNotif.ServiceName)
			elog.Info().Println(value)
		}
	}
}

// ('queue:notification', '{"regnum":"аю88092213","operatorRegnum":"","date":"2023-06-15 16:39:52","serviceName":"WS100101_getCitizenIDCardInfo","serviceDesc":"Иргэний үнэмлэхний мэдээлэл дамжуулах сервис","orgName":"Сангийн яам","requestId":"a43b49c3-ba1e-4eb0-883e-bf04b2ecacfc","resultCode":0,"resultMessage":"амжилттай","clientId":126,"retryCount":0,"civilId":"650187850617"}')
