package routes

import (
	"fmt"
	"net/http"
)

var (
	version   string // app version
	buildTime string // build time
	commitID  string // commit ID
)

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "App Version: %s\n", version)
	fmt.Fprintf(w, "Build Time: %s\n", buildTime)
	fmt.Fprintf(w, "Commit ID: %s\n", commitID)
}
