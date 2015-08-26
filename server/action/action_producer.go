package action

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/crask/mqproxy/global"
	"github.com/crask/mqproxy/producer/kafka"
	"github.com/crask/mqproxy/serializer"

	"github.com/golang/glog"
	//"gopkg.in/Shopify/sarama.v1"
)

func HttpProducerAction(w http.ResponseWriter, r *http.Request) {

	var tspro, tepro, tsall, teall time.Time
	var topic, partitionKey, logId, contentType string
	var resData producer.Response

	tsall = time.Now()
	defer func() {
		teall = time.Now()
		glog.Infof("[kafkaproxy][logid:%s][topic:%s][partition-key:%s][partition:%d][offset:%d][cost:%v][cost_mq:%v]",logId, topic, partitionKey, resData.Data.Partition, resData.Data.Offset, fmt.Sprintf("%.2f", teall.Sub(tsall).Seconds()*1000), fmt.Sprintf("%.2f", tepro.Sub(tspro).Seconds()*1000))
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
		glog.Errorf("[kafkaproxy] Invalid request from %s, topic not found.", r.RemoteAddr)
		echo2client(w, s, producer.Response{
			Errno:  -1,
			Errmsg: "Not Found TOPIC Header",
		})
		return
	}

	partitionKey = r.Header.Get("X-Kmq-Partition-Key")
	if partitionKey == "" {
		glog.Errorf("[kafkaproxy] Invalid request from %s, partition-key not found.", r.RemoteAddr)
		echo2client(w, s, producer.Response{
			Errno:  -1,
			Errmsg: "Not Found PARTITION_KEY Header",
		})
		return
	}

	logId = r.Header.Get("X-Kmq-Logid")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Errorf("[kafkaproxy] Invalid request from %s, read body error[%s].", r.RemoteAddr, err.Error())
		return
	}

	contentType = r.Header.Get("content-type")

	tspro = time.Now()
	prod, i := global.ProducerPool.GetProducer()
	resData, err = prod.SendMessage(producer.Request{
		Topic:        topic,
		PartitionKey: partitionKey,
		TimeStamp:    time.Now().UnixNano() / 1000000,
		Data:         string(body),
		LogId:        logId,
		ContentType:  contentType,
	})
	tepro = time.Now()

	// flush response first
	echo2client(w, s, resData)

	// rebuild connection when error
	if err != nil {
		glog.Errorf("[kafkaproxy][logid:%s][topic:%s][partition-key:%s][Pool-id:%d]Produce Message error, err=%s",
			logId, topic, partitionKey, i, err.Error())
		if err = global.ProducerPool.ReopenProducer(prod, i); err != nil {
			glog.Errorf("[kafkaproxy]Reopen producer failed when failover, pool_id=%d err=%s", i, err.Error())
		}

	}
}

func echo2client(w http.ResponseWriter, s serializer.Serializer, res producer.Response) {
	b, e := s.Marshal(map[string]interface{}{
		"errno":  res.Errno,
		"errmsg": res.Errmsg,
		"data":   res.Data,
	})
	if e != nil {
		glog.Errorf("[kafkaproxy]Marshal http response error, %v", e)
	} else {
		io.WriteString(w, string(b))
	}

	return
}
