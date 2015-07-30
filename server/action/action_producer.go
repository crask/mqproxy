package action

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/crask/mqproxy/global"
	"github.com/crask/mqproxy/producer/kafka"
	"github.com/crask/mqproxy/serializer"

	"github.com/golang/glog"
	"gopkg.in/Shopify/sarama.v1"
)

func HttpProducerAction(w http.ResponseWriter, r *http.Request) {

	var tspro, tepro, tsall, teall time.Time
	var topic, partitionKey, logId string

	tsall = time.Now()
	defer func() {
		teall = time.Now()
		glog.Infof("[kafka-pusher][logid:%s][topic:%s][partition-key:%s][cost:%v][cost_mq:%v]",
			logId, topic, partitionKey, teall.Sub(tsall), tepro.Sub(tspro))
	}()

	// Pasre URL
	query := r.URL.Query()
	msgConverter := strings.Join(query["format"], "")
	if msgConverter == "" {
		msgConverter = "json"
	}

	s := serializer.Serializer{Converter: msgConverter}

	// Get Header
	topic = r.Header.Get("X-Kmq-Topic")
	if topic == "" {
		glog.Errorf("[kafkapusher] Invalid request from %s, topic not found.", r.RemoteAddr)
		echo2client(w, s, producer.Response{
			Errno:  -1,
			Errmsg: "Not Found TOPIC Header",
		})
		return
	}

	partitionKey = r.Header.Get("X-Kmq-Partition-Key")
	if partitionKey == "" {
		glog.Errorf("[kafkapusher] Invalid request from %s, partition-key not found.", r.RemoteAddr)
		echo2client(w, s, producer.Response{
			Errno:  -1,
			Errmsg: "Not Found PARTITION_KEY Header",
		})
		return
	}

	logId = r.Header.Get("X-Kmq-Logid")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Errorf("[kafkapusher] Invalid request from %s, read body error[%s].", r.RemoteAddr, err.Error())
		return
	}

	var resData producer.Response

	defer func() {
		if err := recover(); err != nil {
			glog.Fatalf("catch a panic: %v, recover it", err)
			err = global.ProducerPool.Rebuild()
		}
	}()

	tspro = time.Now()
	resData, err = global.ProducerPool.GetProducer().SendMessage(producer.Request{
		Topic:        topic,
		PartitionKey: partitionKey,
		TimeStamp:    time.Now().UnixNano() / 1000000,
		Data:         string(body),
		LogId:        logId,
	})
	tepro = time.Now()

	switch err.(type) {
	case nil:
		break
	case sarama.PacketEncodingError:
		log.Printf("producer SendMessage error, %v", err)
		break
	default:
		if err == io.EOF {
			panic("producer detected closed LAN connection, panic")
		} else {
			panic(err)
		}
		break
	}

	echo2client(w, s, resData)
}

func echo2client(w http.ResponseWriter, s serializer.Serializer, res producer.Response) {
	b, e := s.Marshal(map[string]interface{}{
		"errno":  res.Errno,
		"errmsg": res.Errmsg,
		"data":   res.Data,
	})
	if e != nil {
		log.Printf("marshal http response error, %v", e)
	} else {
		io.WriteString(w, string(b))
	}

	return
}
