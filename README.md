# mqproxy

###安装
go get github.com/crask/mqproxy

###启动
./mqproxy -c ./proxy.cfg 

###生产消息

```php
<?php
    $ch =curl_init("http://127.0.0.1:9090/produce?format=json");
    curl_setopt($ch, CURLOPT_HEADER, 0);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_BINARYTRANSFER, true);

    $data = array(
        'topic' => 'test',
        'partitionKey' => '123',
        'data' => array(
            'uid'     => 123,
            'uname'   => 'crask',
            'content' => 'welcome to crask',
        ),
    );

    $json_str = json_encode($data);
    curl_setopt($ch, CURLOPT_POSTFIELDS, $json_str);
    $res = curl_exec($ch);
    var_dump($res);
?>
```

###配置项

```shell
#http server(proxy server)配置
[server]
port=9090        
rtimeout=500 #读超时，单位: ms
wtimeout=500 #写超时，单位: ms
maxHeaderBytes=1048576   #http header最大容量，单位: byte
keepAlive=off  # http server是否keep alive

[zookeeper]
addr=127.0.0.1:2181   # kafka依赖的zk server
ctimeout=1   #zk连接超时，单位：sec

[producer]
ctimeout=50  #proxy和kafka broker的连接超时，单位: ms
rtimeout=500 #proxy和kafka broker的读超时，单位: ms
wtimeout=500 #proxy和kafka broker的写超时，单位: ms
keepAlive=10 #proxy和kafka broker连接保持周琦，单位: ms
maxOpenRequest=5 #How many outstanding requests a connection is allowed to have before sending on it blocks (default 5).
partitioner=Hash #producer的均衡策略，支持Hash，Random, RoundRobin
ackStrategy=WaitForLocal  #数据一致策略, NoRespond: 不做任何保证  WaitForLocal: 保证数据提交到leader   WaitForAll: 保证leader和follower强一致
watAckTimeout=500  #ms The maximum duration the broker will wait the receipt of the number of RequiredAcks(default: 1000)
compress=None  #  数据压缩策略，None: 不压缩  Gzip: Gzip压缩   Snappy: Snappy压缩
maxMessageBytes=1000000 # 支持的最大message， 单位: byte
channelBufferSize=0     # 异步提交时采用到，这里写0
poolSize=1   #producer 实例个数，1个就可以了
```
