package server

import (
	"github.com/janqii/config"
	"time"
)

type ProxyConfig struct {
	HttpServerPort           string
	HttpServerReadTimeout    time.Duration
	HttpServerWriteTimeout   time.Duration
	HttpServerMaxHeaderBytes int
	HttpKeepAliveEnabled     bool

	ZookeeperAddr    string
	ZookeeperTimeout time.Duration

	//Producer Configure
	MaxOpenRequests     int
	NetDialTimeout      time.Duration
	NetReadTimeout      time.Duration
	NetWriteTimeout     time.Duration
	NetKeepAlive        time.Duration
	PartitionerStrategy string //Hash, Random, RoundRobin
	WaitAckStrategy     string // NoRespond, WaitForLocal, WaitForAll
	WaitAckTimeoutMs    time.Duration
	CompressionStrategy string //None, Gzip, Snappy
	MaxMessageBytes     int
	ChannelBufferSize   int

	ProducerPoolSize int
}

func NewProxyConfig(file string) (*ProxyConfig, error) {
	cfg := new(ProxyConfig)
	if err := cfg.Parse(file); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *ProxyConfig) Parse(file string) error {
	var err error

	c, err := config.ReadDefault(file)
	if err != nil {
		return err
	}

	// http server cfg
	svrSection := "server"
	cfg.HttpServerPort, err = c.String(svrSection, "port")
	if err != nil {
		return err
	}
	rtimeout, err := c.Int(svrSection, "rtimeout")
	if err != nil {
		rtimeout = 5000
	}
	cfg.HttpServerReadTimeout = time.Duration(rtimeout)

	wtimeout, err := c.Int(svrSection, "wtimeout")
	if err != nil {
		wtimeout = 5000
	}
	cfg.HttpServerWriteTimeout = time.Duration(wtimeout)

	cfg.HttpServerMaxHeaderBytes, err = c.Int(svrSection, "maxHeaderBytes")
	if err != nil {
		cfg.HttpServerMaxHeaderBytes = 1048576
	}
	cfg.HttpKeepAliveEnabled, err = c.Bool(svrSection, "keepAlive")
	if err != nil {
		cfg.HttpKeepAliveEnabled = false
	}

	//zookeeper cfg
	zkSection := "zookeeper"
	cfg.ZookeeperAddr, err = c.String(zkSection, "addr")
	if err != nil {
		return err
	}
	zkctimeout, err := c.Int(zkSection, "ctimeout")
	if err != nil {
		zkctimeout = 1
	}
	cfg.ZookeeperTimeout = time.Duration(zkctimeout) * time.Second

	//producer cfg
	pdSection := "producer"
	cfg.MaxOpenRequests, err = c.Int(pdSection, "maxOpenRequests")
	if err != nil {
		cfg.MaxOpenRequests = 5
	}

	dialTimeout, err := c.Int(pdSection, "ctimeout")
	if err != nil {
		dialTimeout = 500
	}
	cfg.NetDialTimeout = time.Duration(dialTimeout)

	netrtimout, err := c.Int(pdSection, "rtimeout")
	if err != nil {
		netrtimout = 500
	}
	cfg.NetReadTimeout = time.Duration(netrtimout)

	netwtimeout, err := c.Int(pdSection, "wtimeout")
	if err != nil {
		netwtimeout = 500
	}
	cfg.NetWriteTimeout = time.Duration(netwtimeout)

	netKeepAlive, err := c.Int(pdSection, "keepAlive")
	if err != nil {
		netKeepAlive = 0
	}
	cfg.NetKeepAlive = time.Duration(netKeepAlive)

	cfg.PartitionerStrategy, err = c.String(pdSection, "partitioner")
	if err != nil {
		cfg.PartitionerStrategy = "Hash"
	}
	cfg.WaitAckStrategy, err = c.String(pdSection, "ackStrategy")
	if err != nil {
		cfg.WaitAckStrategy = "WaitForLocal"
	}
	acktimeout, err := c.Int(pdSection, "waitAckTimeout")
	if err != nil {
		acktimeout = 1000
	}
	cfg.WaitAckTimeoutMs = time.Duration(acktimeout)
	cfg.CompressionStrategy, err = c.String(pdSection, "compress")
	if err != nil {
		cfg.CompressionStrategy = "None"
	}
	cfg.MaxMessageBytes, err = c.Int(pdSection, "maxMessageBytes")
	if err != nil {
		cfg.MaxMessageBytes = 1000000
	}
	cfg.ChannelBufferSize, err = c.Int(pdSection, "channelBufferSize")
	if err != nil {
		cfg.ChannelBufferSize = 0
	}
	cfg.ProducerPoolSize, err = c.Int(pdSection, "poolSize")
	if err != nil {
		cfg.ProducerPoolSize = 1
	}

	return nil
}
