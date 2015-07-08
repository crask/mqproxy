package action

import (
	"io"
	"net/http"
	//"io/ioutil"
)

func PartitionsAction(w http.ResponseWriter, r *http.Request) {
	//TODO:
	io.WriteString(w, "PartitionsAction!")
}
