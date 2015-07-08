package action

import (
	"io"
	"net/http"
	//"io/ioutil"
)

func StatsAction(w http.ResponseWriter, r *http.Request) {
	//TODO:
	io.WriteString(w, "StatsAction!")
}
