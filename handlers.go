package main

import (
	"fmt"
	"net/http"
)

func handleGetProxyList(w http.ResponseWriter, r *http.Request) {
	returnResult(w, proxyList.Proxies)
}

func handleGetProxyInfo(w http.ResponseWriter, r *http.Request) {
	proxyName := r.URL.Query().Get("name")

	if proxyName == "" {
		returnError(w, 400, "Bad Request", "Querystring value for parameter 'name' not provided")
		return
	}

	p, err := proxyList.GetProxyByName(proxyName)
	if err != nil {
		returnError(w, 404, "Not Found", err.Error())
		return
	}

	returnResult(w, p)
}

func handleGetProxyStart(w http.ResponseWriter, r *http.Request) {
	proxyName := r.URL.Query().Get("name")

	if proxyName == "" {
		returnError(w, 400, "Bad Request", "Querystring value for parameter 'name' not provided")
		return
	}

	p, err := proxyList.GetProxyByName(proxyName)
	if err != nil {
		returnError(w, 404, "Not Found", fmt.Sprintf("Proxy '%s' not found", proxyName))
		return
	}

	if p.Active {
		returnError(w, 400, "Bad Request", fmt.Sprintf("Proxy '%s' already started", p.Name))
		return
	}

	if err := p.Start(); err != nil {
		returnError(w, 500, "Internal Server Error", fmt.Sprintf("Proxy '%s' could not be started: %s", p.Name, err.Error()))
		return
	}
	returnResult(w, p)
}

func handleGetProxyStop(w http.ResponseWriter, r *http.Request) {
	proxyName := r.URL.Query().Get("name")

	if proxyName == "" {
		returnError(w, 400, "Bad Request", "Querystring value for parameter 'name' not provided")
		return
	}

	p, err := proxyList.GetProxyByName(proxyName)
	if err != nil {
		returnError(w, 404, "Not Found", fmt.Sprintf("Proxy '%s' not found", proxyName))
		return
	}

	if !p.Active {
		returnError(w, 400, "Bad Request", fmt.Sprintf("Proxy '%s' already stopped", p.Name))
		return
	}

	if err := p.Stop(); err != nil {
		returnError(w, 500, "Internal Server Error", fmt.Sprintf("Proxy '%s' could not be stopped: %s", p.Name, err.Error()))
		return
	}
	returnResult(w, p)
}
