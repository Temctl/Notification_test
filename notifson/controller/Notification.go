package controller

import (
	"fmt"
	"net/http"
)

func Notification(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Fdsfsd")
}
