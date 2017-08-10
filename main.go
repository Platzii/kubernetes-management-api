package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/Sirupsen/logrus"
	"github.com/koding/multiconfig"
	"github.com/rs/cors"
)

// Config contains the overall configuration
type Config struct {
	KubeConfig         string
	KubeCtl            string
	Port               string `default:"8000"`
	kubeConfigLocation string
	kubeCtlLocation    string
}

var config *Config
var proxyList ProxyList

func main() {
	mc := multiconfig.New()

	config = &Config{}
	mc.MustLoad(config)

	// get kube config location
	if config.KubeConfig != "" {
		config.kubeConfigLocation = config.KubeConfig
	} else {
		if _, err := os.Stat("kubeconfig"); !os.IsNotExist(err) {
			config.kubeConfigLocation = "kubeconfig"
		} else {
			usr, err := user.Current()
			if err != nil {
				logrus.Fatalf("Could not get current user info: %s", err.Error())
			}
			path := filepath.Join(usr.HomeDir, ".kube", "config")
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				config.kubeConfigLocation = path
			} else {
				logrus.Fatalln("FATAL ERROR\n\r\n\rCould not find Kube config file\n\rIt should be specified with --kubeconfig OR be located at ./kubeconfig or ~/.kube/config\n\r")
			}
		}
	}
	logrus.Infof("Using Kubernetes config at %s", config.kubeConfigLocation)

	// get kubectl location
	if config.KubeCtl != "" {
		config.kubeCtlLocation = config.KubeCtl
	} else {
		var err error
		config.kubeCtlLocation, err = exec.LookPath("kubectl")
		if err != nil {
			logrus.Fatalf("Could not find kubectl binary")
			logrus.Fatalln("FATAL ERROR\n\r\n\rCould not find kubectl executable\n\rIt should be discoverable using $PATH OR be specified with --kubectl\n\r")
		}
	}
	logrus.Infof("Using kubectl binary at %s", config.kubeCtlLocation)

	if err := proxyList.FillProxies(); err != nil {
		logrus.Fatalf("Could not get proxy list: %s", err.Error())
	}
	logrus.Infof("Found %d contexts", len(proxyList.Proxies))

	mux := http.NewServeMux()
	mux.HandleFunc("/proxy/list", handleGetProxyList)
	mux.HandleFunc("/proxy/info", handleGetProxyInfo)
	mux.HandleFunc("/proxy/start", handleGetProxyStart)
	mux.HandleFunc("/proxy/stop", handleGetProxyStop)

	handler := cors.Default().Handler(mux)
	logrus.Infof("Starting to listen on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, handler); err != nil {
		logrus.Fatalf("Could not start HTTP server: %s", err.Error())
	}
}

func getContexts() ([]string, error) {
	kubeConfig, err := clientcmd.LoadFromFile(config.kubeConfigLocation)
	if err != nil {
		return nil, fmt.Errorf("load kubeconfig: %s", err.Error())
	}

	var contexts []string
	for contextName := range kubeConfig.Contexts {
		contexts = append(contexts, contextName)
	}

	sort.Strings(contexts)

	return contexts, nil
}
