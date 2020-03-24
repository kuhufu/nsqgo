nsqgo
------------------------------
基于[github.com/nsqio/go-nsq](github.com/nsqio/go-nsq)

基于nsqlookupd提供自动发现nsqd的能力。

对多个nsqd负载均衡：随机，轮转。

```go
import github.com/kuhufu/nsqgo
```

```go
producer, err := nsqgo.NewProducer(
    nsqgo.WithNSQLookupd("localhost:4161"),
    nsqgo.WithNSQConfig(nsq.NewConfig()),
)

err = producer.Publish("fruit", "apple")
if err != nil {
    log.Println(err)
}
```





