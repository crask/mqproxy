package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/crask/mqproxy/internal/version"
	"github.com/crask/mqproxy/server"

	"github.com/golang/glog"
)

var (
	configFile string
	vers       bool
	help       bool
)

func showVersion() {
	fmt.Println(version.String("kafkaproxy"))
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
