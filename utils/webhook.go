package utils

import (
	"fmt"
	"net/http"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
