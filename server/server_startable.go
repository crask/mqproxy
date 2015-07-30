package server

import (
	"net/http"
	"runtime"
	"sync"

	"github.com/crask/mqproxy/global"
	"github.com/crask/mqproxy/producer/kafka"
	"github.com/crask/mqproxy/server/router"

	"github.com/golang/glog"
	"github.com/wvanbergen/kazoo-go"
)

func Startable(cfg *ProxyConfig) error {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	wg := new(sync.WaitGroup)

	zkNodes, chroot := kazoo.ParseConnectionString(cfg.ZookeeperAddr)
	zkConfig := &kazoo.Config{
		Chroot:  chroot,
		Timeout: cfg.ZookeeperTimeout,
	}

	zkProxy, err := kazoo.NewKazoo(zkNodes, zkConfig)
	if err != nil {
		glog.Errorf("[kafkaproxy]NewKazoo error: %v", err)
		return err
	}
	defer zkProxy.Close()

	brokerList, err := zkProxy.BrokerList()
	if err != nil {
		glog.Errorf("[kafkaproxy]Get broker list error: %v", err)
		return err
	}

	pcfg := &producer.KafkaProducerConfig{
		Addrs:               brokerList,
		MaxOpenRequests:     cfg.MaxOpenRequests,
		DialTimeout:         cfg.NetDialTimeout,
		ReadTimeout:         cfg.NetReadTimeout,
		WriteTimeout:        cfg.NetWriteTimeout,
		KeepAlive:           cfg.NetKeepAlive,
		PartitionerStrategy: cfg.PartitionerStrategy,
		WaitAckStrategy:     cfg.WaitAckStrategy,
		WaitAckTimeoutMs:    cfg.WaitAckTimeoutMs,
		CompressionStrategy: cfg.CompressionStrategy,
		MaxMessageBytes:     cfg.MaxMessageBytes,
		ChannelBufferSize:   cfg.ChannelBufferSize,
	}

	global.ProducerPool, err = global.NewKafkaProducerPool(pcfg, cfg.ProducerPoolSize)
	defer global.DestoryKafkaProducerPool(global.ProducerPool)

	if err != nil {
		glog.Errorf("[kafkaproxy]Create kafka producer pool error: %v", err)
		return err
	}

	var (
		//statMux  map[string]func(http.ResponseWriter, *http.Request)
		proxyMux map[string]func(http.ResponseWriter, *http.Request)
	)
	//statMux = make(map[string]func(http.ResponseWriter, *http.Request))
	proxyMux = make(map[string]func(http.ResponseWriter, *http.Request))

	proxyHttpServer := &HttpServer{
		Addr:            ":" + cfg.HttpServerPort,
		Handler:         &HttpHandler{Mux: proxyMux},
		ReadTimeout:     cfg.HttpServerReadTimeout,
		WriteTimeout:    cfg.HttpServerWriteTimeout,
		MaxHeaderBytes:  cfg.HttpServerMaxHeaderBytes,
		KeepAliveEnable: cfg.HttpKeepAliveEnabled,
		RouterFunc:      router.ProxyServerRouter,
		Wg:              wg,
		Mux:             proxyMux,
	}

	proxyHttpServer.Startup()
	defer proxyHttpServer.Shutdown()

	glog.V(2).Info("[tracing]MQ Proxy is running...")
	wg.Wait()
	glog.V(2).Info("[tracing]MQ Proxy is exiting...")

	return nil
}
