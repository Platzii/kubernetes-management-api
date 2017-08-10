package main

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/Sirupsen/logrus"
)

// Proxy represents a proxy connection
type Proxy struct {
	Name   string `json:"name"`
	Port   string `json:"port"`
	Active bool   `json:"active"`
	PID    int    `json:"pid"`
	cmd    *exec.Cmd
}

// Start Starts a proxy connection
func (p *Proxy) Start() error {
	logrus.Debugf("Start proxy for context '%s'", p.Name)

	p.cmd = exec.Command(config.kubeCtlLocation, "--context", p.Name, "proxy", "-p", p.Port)
	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("start: %s", err.Error())
	}

	p.Active = true
	p.PID = p.cmd.Process.Pid

	logrus.Debugf("Proxy process for context '%s' running with PID '%d'", p.Name, p.PID)

	return nil
}

// Stop Stops a proxy connection
func (p *Proxy) Stop() error {
	logrus.Debugf("Stop proxy for context '%s'", p.Name)

	// if err := p.cmd.Process.Signal(os.Interrupt); err != nil { // SIGINT not working on Windows
	if err := p.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("kill process: %s", err.Error())
	}
	p.Active = false
	p.PID = -1
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// ProxyList represents a list of proxy connections
type ProxyList struct {
	Proxies []*Proxy
}

// FillProxies Fills Proxies with Proxy objects received from entries in kubeconfig file
func (pl *ProxyList) FillProxies() error {
	contexts, err := getContexts()
	if err != nil {
		return fmt.Errorf("get contexts: %s", err.Error())
	}

	for i, context := range contexts {
		p := &Proxy{
			Name:   context,
			Port:   strconv.Itoa(8001 + i),
			Active: false, // TODO: auto detect if proxy is already running?
			PID:    -1,
		}
		pl.Proxies = append(pl.Proxies, p)
	}

	return nil
}

// GetProxyByName Returns Proxy object for given name
func (pl *ProxyList) GetProxyByName(name string) (*Proxy, error) {
	for _, proxy := range proxyList.Proxies {
		if proxy.Name == name {
			return proxy, nil
		}
	}
	return nil, fmt.Errorf("Could not find proxy with name %s", name)
}
