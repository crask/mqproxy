package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/crask/mqproxy/server"

	"github.com/golang/glog"
)

const (
	VERSION = "1.0.3"
)

var (
	configFile string
	vers       bool
	help       bool
	gitCommit  string
)

func showVersion() {
	fmt.Println(fmt.Sprintf("%s-%s", VERSION, gitCommit))
	flag.Usage()
}

func init() {
	flag.StringVar(&configFile, "c", "conf/proxy.cfg", "Set config file path")
	flag.BoolVar(&vers, "V", false, "Show version")
	flag.BoolVar(&help, "h", false, "Show help")
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if vers || help {
		showVersion()
		os.Exit(0)
	}

	glog.Info("[kafkaproxy]Server starting...")

	cfg, err := server.NewProxyConfig(configFile)
	if err != nil {
		glog.Errorf("[kafkaproxy]parse config error, %v", err)
		os.Exit(0)
	}

	if err = server.Startable(cfg); err != nil {
		glog.Errorf("[kafkaproxy]server startable error, %v", err)
		os.Exit(0)
	}
}
