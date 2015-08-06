package action

import (
	"net/http"
	//"io/ioutil"
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/toolbox"
)

func ProfileAction(w http.ResponseWriter, r *http.Request) {
	//TODO:
	r.ParseForm()
	command := r.Form.Get("command")
	format := r.Form.Get("format")
	data := make(map[string]interface{})
	var result bytes.Buffer
	if command != "" {
		toolbox.ProcessInput(command, &result)
		data["Content"] = result.String()

		if format == "json" {
			dataJson, err := json.Marshal(data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(dataJson)
			return
		}
	}
}
