package action

import (
	"github.com/crask/mqproxy/global"
	"github.com/crask/mqproxy/producer/kafka"
	"github.com/crask/mqproxy/serializer"
	"gopkg.in/Shopify/sarama.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func HttpProducerAction(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read request body error")
		return
	}

	query := r.URL.Query()
	msgConverter := strings.Join(query["format"], "")
	if msgConverter == "" {
		msgConverter = "json"
	}

	s := serializer.Serializer{Converter: msgConverter}

	var resData producer.Response
	var reqData producer.Request
	if err = s.Unmarshal(body, &reqData); err != nil {
		log.Printf("Unmarshal HttpRequest error, %v", err)
		resData = producer.Response{
			-1,
			"unmarshal http request data error",
			make([]producer.MessageLocation, 0),
		}
	} else {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("catch a panic: %v, recover it", err)
				err = global.ProducerPool.Rebuild()
				//				echo2client(w, s, resData, err)
			}
		}()
		resData, err = global.ProducerPool.GetProducer().SendMessage(reqData)
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
	}

	echo2client(w, s, resData, err)
}

func echo2client(w http.ResponseWriter, s serializer.Serializer, res producer.Response, e error) {
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
