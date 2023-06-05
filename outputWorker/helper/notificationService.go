package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Temctl/E-Notification/outputWorker/helper"
	"github.com/Temctl/E-Notification/util/model"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func getPushToken(regnum string) string {

	return "dIMtXp4UUkdZoj1D4M8wwD:APA91bFzD_WEW2cvd6QaXRk9cllEbr_ECrREZ2KzlbjbbWpW-7I5gNYgpgZOLGUu4HpNtc_hjyPG6YYceUbjhniqQmafV-DXV5__ezlMo07-Wq1m0trdJ5H7UWPe9SgxeFmjwN8HwmBO"
}

func send