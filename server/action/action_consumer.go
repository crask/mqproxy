package action

import (
	"io"
	"net/http"
	//"io/ioutil"
)

func FetchMessageAction(w http.ResponseWriter, r *http.Request) {
	//TODO:
	io.WriteString(w, "FetchMessageAction!")
}
