package worker

import (
	"encoding/json"
	"fmt"
	"log"

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

	// ----------------------------------------------------------------------
	// DB CONNECTION --------------------------------------------------------
	// ----------------------------------------------------------------------
	db, err := connections.ConnectPostgreSQL()
	if err != nil {
		elog.Error().Fatal(err)
	}
	stmt, err := db.Prepare("INSERT INTO xypnotification (regnum, operatorregnum,date, servicename, servicedesc,orgname, requestid, resultcode, civilid, clientid) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
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
				elog.Error().Println("Error:", err)
				continue
			}
			_, err = stmt.Exec(
				xypNotif.Regnum,
				xypNotif.OperatorRegnum,
				xypNotif.Date,
				xypNotif.ServiceName,
				xypNotif.ServiceDesc,
				xypNotif.OrgName,
				xypNotif.RequestId,
				xypNotif.ResultCode,
				xypNotif.CivilId,
				xypNotif.ClientId)
			if err != nil {
				elog.Error().Println("Error:", err)
				continue
			}

			fmt.Println(xypNotif)
		}
	}
}

// ('queue:notification', '{"regnum":"уц03211011","operatorRegnum":"","date":"2023-06-15 16:39:52","serviceName":"WS100101_getCitizenIDCardInfo","serviceDesc":"Иргэний үнэмлэхний мэдээлэл дамжуулах сервис","orgName":"Сангийн яам","requestId":"a43b49c3-ba1e-4eb0-883e-bf04b2ecacfc","resultCode":0,"resultMessage":"амжилттай","clientId":126,"retryCount":0,"civilId":"891834062934"}')
