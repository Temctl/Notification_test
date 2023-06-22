package sendworker

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/streadway/amqp"
)



func init(){
	file, _ := os.Create("./log/output.log")
	log.SetOutput(file)
	file.Close()
	log.SetFlags(log.Ldate | log.Lshortfile)

}


func main(){

	redis, err := connections.ConnectionRedis()
	if err != nil{
		//log error
	}
	//get mongodb connection
	mongoDB, err := connections.ConnectMongoDB()
	if err != nil{
		//log
	}


	var wg sync.WaitGroup
	wg.Add(3)
	// ----------------------------------------------------------------------
	// CRON JOB -------------------------------------------------------------
	// ----------------------------------------------------------------------
	go func() {
		defer wg.Done()
		worker.AttentionNotificationEveryday()
	}()
	// ----------------------------------------------------------------------
	// XYP WORKER -----------------------------------------------------------
	// ----------------------------------------------------------------------
	go func() {
		defer wg.Done()
		XypFromDb()
	}()

	// Wait for all goroutines to complete
	wg.Wait()
}