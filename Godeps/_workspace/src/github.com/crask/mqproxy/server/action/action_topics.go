package action

import (
	"io"
	"net/http"
	//"io/ioutil"
)

func TopicsAction(w http.ResponseWriter, r *http.Request) {
	//TODO:
	io.WriteString(w, "TopicsAction!")
}
