package router

import (
	"github.com/crask/mqproxy/server/action"
	"net/http"
)

func StatServerRouter(mux map[string]func(http.ResponseWriter, *http.Request)) {
	mux["/stats"] = action.StatsAction
}
