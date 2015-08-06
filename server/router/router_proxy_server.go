package router

import (
	"github.com/crask/mqproxy/server/action"
	"net/http"
)

func ProxyServerRouter(mux map[string]func(http.ResponseWriter, *http.Request)) {
	mux["/topics"] = action.TopicsAction

	mux["/produce"] = action.HttpProducerAction
	mux["/consumer/fetch"] = action.FetchMessageAction
	mux["/profile"] = action.ProfileAction
}
