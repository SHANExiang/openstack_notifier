package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"sincerecloud.com/openstack_notifier/consts"
	"sincerecloud.com/openstack_notifier/global"
	"sincerecloud.com/openstack_notifier/nonick"
	"sincerecloud.com/openstack_notifier/services/consume"
	"sincerecloud.com/openstack_notifier/services/publish"
	"sincerecloud.com/openstack_notifier/utils"
	"strings"
	"time"
)

var PublisherMap = make(map[string]*publish.Publisher)
var Consumers = make(map[string]*consume.RabbitMQ)
var OpenstackMQs = make(map[string]nonick.MQInfo)


func GetMQInfos() (map[string]nonick.MQInfo, error) {
	var ret = make(map[string]nonick.MQInfo)
	getUrl := fmt.Sprintf("http://%s%s",
		global.CONF.ECS.NovaURL, consts.GetMqInfoSuffix)
	global.LOG.Info(fmt.Sprintf("GetMQInfos request url (%+v)", getUrl))
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(getUrl)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		global.LOG.Error(fmt.Sprintf("GetMQInfos err %+v", err))
		return map[string]nonick.MQInfo{}, err
	}
	var mqInfoResp nonick.MqInfoResp
	content := resp.Body()
	global.LOG.Debug(fmt.Sprintf("GetMQInfos response content %+v", string(content)))
	if err := json.Unmarshal(content, &mqInfoResp); err != nil {
		global.LOG.Error(fmt.Sprintf("Parse json err %+v", err))
		return map[string]nonick.MQInfo{}, err
	}
	global.LOG.Info(fmt.Sprintf("GetMQInfos response body %+v", mqInfoResp))
	for _, info := range mqInfoResp.Data {
		key := info.AzCode
		if len(key) != 0 {
			ret[key] = info
		}
	}
	return ret, nil
}


func InitPublisher(resourceType string) *publish.Publisher {
	var url string
	switch resourceType {
	case consts.NOVA:
	    url = global.CONF.ECS.NovaURL + consts.NovaSuffix
	case consts.CINDER:
		url = global.CONF.ECS.CinderURL + consts.CinderSuffix
	case consts.NEUTRON:
		url = global.CONF.ECS.NeutronURL + consts.NeutronSuffix
	case consts.COMPUTE:
		url = global.CONF.ECS.ComputeURL + consts.ComputeSuffix
	}
	publisher := &publish.Publisher{
		Url: "http://" + url,
	}
	global.LOG.Info(fmt.Sprintf("Init publisher %+v success", resourceType))
	return publisher
}

func constructAMQPURI(mqInfo nonick.MQInfo) []string {
	var amqpURIPool = make([]string, 0)
	for _, uri := range strings.Split(strings.TrimSpace(mqInfo.MqEndpoint), ",") {
		amqpURI := fmt.Sprintf("amqp://%s:%s@%s/", mqInfo.MqUserName, mqInfo.MqPassword, strings.TrimSpace(uri))
		amqpURIPool = append(amqpURIPool, amqpURI)
	}
	return amqpURIPool
}



func InitConsumer(mqInfo nonick.MQInfo) *consume.RabbitMQ {
    amqpURIPool := constructAMQPURI(mqInfo)
	consumer := consume.NewRabbitMQ(
		mqInfo.MqChannel,
		mqInfo.MqChannel,
		"#",
		amqpURIPool)
	consumer.AzCode = mqInfo.AzCode
	consumer.Tag = mqInfo.MqChannel + "_" + mqInfo.AzCode
	consumer.Running = true
	consumer.IsSDN = mqInfo.SDN
	global.LOG.Info(fmt.Sprintf("Init consumer %+v success", consumer.Tag))
	return consumer
}

func startConsumer(consumer *consume.RabbitMQ) {
	for {
		global.LOG.Info("start to consume", zap.String("tag", consumer.Tag))
		consume.Consume(consumer, PublisherMap)
		if !consumer.Running {
			break
		}
		time.Sleep(3 * time.Second)
	}
}

func retryPublish(buffer chan *publish.PublishMsg) {
	for {
        msg, ok := <-buffer
        if ok {
            publisher := PublisherMap[consume.GetResourceMap(msg.ResourceType)]
            body, _ := json.Marshal(msg)
            var err = errors.New("publish failed")
            reqID := utils.GenerateRequestID()
            ctx := utils.WrapRequestID(reqID)
            for err != nil {
                if err = publish.RetryDo(func(ctx context.Context, body []byte) error {
                    return publisher.Publish(ctx, body)
                }, ctx, body, publish.DefaultRetries, publish.DefaultSleep); err != nil {
                    time.Sleep(5*time.Second)
                }
            }
		} else {
			global.LOG.Info("exit retry publish channel")
			break
		 }
	}
}

func validateAMQPURI(mqInfo nonick.MQInfo) bool {
	if mqInfo.MqEndpoint == "" || mqInfo.MqUserName == "" || mqInfo.MqPassword == "" {
		return false
	} else {
		return true
	}
}

func InitChannelsAndRetry() (buffers []chan *publish.PublishMsg) {
	for i := 0;i < 6;i++ {
		buffer := make(chan *publish.PublishMsg, 2<<10)
		buffers = append(buffers, buffer)
		go retryPublish(buffer)
	}
	return buffers
}

func call(mqInfo nonick.MQInfo) {
	consumer := InitConsumer(mqInfo)
	Consumers[mqInfo.AzCode] = consumer
	if validateAMQPURI(mqInfo) {
		buffers := InitChannelsAndRetry()
		consumer.Buffers = buffers
		go startConsumer(consumer)
	} else {
		global.LOG.Info(fmt.Sprintf("consumer %s amqpURI is invalid, don't start to consume", consumer.Tag))
	}
}

func getDiffMap(map1, map2 map[string]nonick.MQInfo) (map[string]nonick.MQInfo, map[string]nonick.MQInfo, map[string]nonick.MQInfo) {
	addMQs := make(map[string]nonick.MQInfo)
	deleteMQs := make(map[string]nonick.MQInfo)
	updateMQs := make(map[string]nonick.MQInfo)

	for key, value := range map1 {
		if v, ok := map2[key]; ok {
			updateMQs[key] = v
		} else {
			deleteMQs[key] = value
		}
	}

	for key, value := range map2 {
		if _, ok := map1[key]; !ok {
			addMQs[key] = value
		}
	}
	return addMQs, updateMQs, deleteMQs
}

func ReloadMQCall() {
	ticker := time.NewTicker(60 * time.Second)
	for {
		<-ticker.C
		mqInfos, err := GetMQInfos()
		if err != nil {
			global.LOG.Error("GetMQInfos failed : %+v", zap.Error(err))
		} else {
			addMQs, updateMQs, deleteMQs := getDiffMap(OpenstackMQs, mqInfos)
			for acCode, mqInfo := range addMQs {
				global.LOG.Info(fmt.Sprintf("add a openstack cluster %s", acCode))
				call(mqInfo)
				OpenstackMQs[acCode] = mqInfo
			}
			for acCode, mqInfo := range updateMQs {
				if mqInfo.MqEndpoint != OpenstackMQs[acCode].MqEndpoint ||
					mqInfo.MqUserName != OpenstackMQs[acCode].MqUserName ||
					mqInfo.MqPassword != OpenstackMQs[acCode].MqPassword ||
					mqInfo.MqChannel != OpenstackMQs[acCode].MqChannel ||
					mqInfo.SDN != OpenstackMQs[acCode].SDN {
					global.LOG.Info(fmt.Sprintf("update a openstack cluster %s", acCode))
					Consumers[acCode].Close()
					time.Sleep(3 * time.Second)
					Consumers[acCode].AmqpURIPool = constructAMQPURI(mqInfo)
					Consumers[acCode].QueueName = mqInfo.MqChannel
					Consumers[acCode].Exchange = mqInfo.MqChannel
					Consumers[acCode].Tag = mqInfo.MqChannel + "_" + mqInfo.AzCode
					Consumers[acCode].IsSDN = mqInfo.SDN
					OpenstackMQs[acCode] = mqInfo
					if validateAMQPURI(mqInfo) {
						Consumers[acCode].Running = true
						if len(Consumers[acCode].Buffers) == 0 {
							Consumers[acCode].Buffers = InitChannelsAndRetry()
						}
						go startConsumer(Consumers[acCode])
					} else {
						global.LOG.Info(fmt.Sprintf("consumer %s amqpURI is invalid, don't start to consume",
							Consumers[acCode].Tag))
					}
				}
			}

			for acCode, _ := range deleteMQs {
				global.LOG.Info(fmt.Sprintf("delete a openstack cluster %s", acCode))
				for _, buffer := range Consumers[acCode].Buffers {
					close(buffer)
				}
				Consumers[acCode].Close()
				delete(OpenstackMQs, acCode)
			}
		}
	}
	ticker.Stop()
}

func RunServer() {
	PublisherMap[consts.NOVA] = InitPublisher(consts.NOVA)
	PublisherMap[consts.CINDER] = InitPublisher(consts.CINDER)
	PublisherMap[consts.NEUTRON] = InitPublisher(consts.NEUTRON)
	PublisherMap[consts.COMPUTE] = InitPublisher(consts.COMPUTE)
	var err error
	OpenstackMQs, err = GetMQInfos()
	if err != nil {
		global.LOG.Error("GetMQInfos failed", zap.Error(err))
	}
	global.LOG.Info(fmt.Sprintf("OpenstackMQs==%+v", OpenstackMQs))
	for _, mqInfo := range OpenstackMQs {
		call(mqInfo)
	}
	go ReloadMQCall()
}
