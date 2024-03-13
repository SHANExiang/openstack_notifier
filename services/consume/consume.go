package consume

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"sincerecloud.com/openstack_notifier/consts"
	"sincerecloud.com/openstack_notifier/global"
	"sincerecloud.com/openstack_notifier/services/publish"
	"sincerecloud.com/openstack_notifier/utils"
	"time"
)

func Consume(consumer *RabbitMQ, publisherMap map[string]*publish.Publisher) {
	var err error
	if err = consumer.GetChannel(); err != nil {
		return
	}
	global.LOG.Info(fmt.Sprintf("consumer %s got channel, starting Consume (consumer tag %q)", consumer.Tag, ""))
	// nova queue nova_ecs_voneyun_topic
	novaQueue := fmt.Sprintf("%s_%s.info", consts.NOVA, consumer.QueueName)
	novaInfoQueue, err := consumer.QueueDeclare(fmt.Sprintf(novaQueue))
	if err != nil {
		global.LOG.Error(fmt.Sprintf("%s %s queue Declare err: %s", consumer.Tag, novaQueue, err))
		return
	}
	global.LOG.Info(fmt.Sprintf("%s declared Queue (%q %d messages, %d consumers)," +
		" binding to Exchange (key %q)",
		consumer.Tag, novaInfoQueue.Name, novaInfoQueue.Messages, novaInfoQueue.Consumers, consumer.Key))
	novaErrQueueName := fmt.Sprintf("%s_%s.error", consts.NOVA, consumer.QueueName)
	novaErrQueue, err := consumer.QueueDeclare(novaErrQueueName)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("%s %s queue Declare err: %s", consumer.Tag, novaErrQueueName, err))
		return
	}
	global.LOG.Info(fmt.Sprintf("%s declared Queue (%q %d messages, %d consumers)," +
		" binding to Exchange (key %q)",
		consumer.Tag, novaErrQueue.Name, novaErrQueue.Messages, novaErrQueue.Consumers, consumer.Key))

	// cinder queue cinder_ecs_voneyun_topic
	cinderInfoQueueName := fmt.Sprintf("%s_%s.info", consts.CINDER, consumer.QueueName)
	cinderInfoQueue, err := consumer.QueueDeclare(fmt.Sprintf(cinderInfoQueueName))
	if err != nil {
		global.LOG.Error(fmt.Sprintf("%s %s queue Declare err: %s", consumer.Tag, cinderInfoQueueName, err))
		return
	}
	global.LOG.Info(fmt.Sprintf("%s declared Queue (%q %d messages, %d consumers)," +
		" binding to Exchange (key %q)",
		consumer.Tag, cinderInfoQueue.Name, cinderInfoQueue.Messages, cinderInfoQueue.Consumers, consumer.Key))
	cinderErrQueueName := fmt.Sprintf("%s_%s.error", consts.CINDER, consumer.QueueName)
	cinderErrQueue, err := consumer.QueueDeclare(cinderErrQueueName)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("%s %s queue Declare err: %s", consumer.Tag, cinderErrQueueName, err))
		return
	}
	global.LOG.Info(fmt.Sprintf("%s declared Queue (%q %d messages, %d consumers)," +
		" binding to Exchange (key %q)",
		consumer.Tag, cinderErrQueue.Name, cinderErrQueue.Messages, cinderErrQueue.Consumers, consumer.Key))

	// neutron queue neutron_ecs_voneyun_topic
	neutronInfoQueueName := fmt.Sprintf("%s_%s.info", consts.NEUTRON, consumer.QueueName)
	neutronInfoQueue, err := consumer.QueueDeclare(fmt.Sprintf(neutronInfoQueueName))
	if err != nil {
		global.LOG.Error(fmt.Sprintf("%s %s queue Declare err: %s", consumer.Tag, neutronInfoQueueName, err))
		return
	}
	global.LOG.Info(fmt.Sprintf("%s declared Queue (%q %d messages, %d consumers)," +
		" binding to Exchange (key %q)",
		consumer.Tag, neutronInfoQueue.Name, neutronInfoQueue.Messages, neutronInfoQueue.Consumers, consumer.Key))
	neutronErrQueueName := fmt.Sprintf("%s_%s.error", consts.NEUTRON, consumer.QueueName)
	neutronErrQueue, err := consumer.QueueDeclare(neutronErrQueueName)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("%s %s queue Declare err: %s", consumer.Tag, neutronErrQueueName, err))
		return
	}
	global.LOG.Info(fmt.Sprintf("%s declared Queue (%q %d messages, %d consumers)," +
		" binding to Exchange (key %q)",
		consumer.Tag, neutronErrQueue.Name, neutronErrQueue.Messages, neutronErrQueue.Consumers, consumer.Key))

	global.LOG.Info(fmt.Sprintf("consumer %s starting Consume (consumer tag %q)", consumer.Tag, ""))
	var queues = []amqp.Queue{novaInfoQueue, novaErrQueue, cinderInfoQueue, cinderErrQueue, neutronInfoQueue, neutronErrQueue}
	for index, queue := range queues {
		consumer.wg.Add(1)
		deliveries, err := consumer.Channel.Consume(
			queue.Name,               // name
			"",             // consumerTag,
			false,           // noAck
			false,          // exclusive
			false,           // noLocal
			false,           // noWait
			nil,               // arguments
		)
		if err != nil {
			global.LOG.Error(fmt.Sprintf("%s consume %s err: %s", consumer.Tag, queue.Name, err))
			consumer.wg.Done()
			return
		}
		go handle(deliveries, consumer, publisherMap, consumer.Buffers[index])
	}
	consumer.wg.Wait()
}

func handle(deliveries <-chan amqp.Delivery, consumer *RabbitMQ, publisherMap map[string]*publish.Publisher, buffer chan *publish.PublishMsg) {
	defer consumer.wg.Done()
	for d := range deliveries {
		reqID := utils.GenerateRequestID()
		global.LOG.Debug(fmt.Sprintf(
			"[req-%s]Consumer %s got %dB delivery: [%v] %q",
			reqID,
			consumer.Tag,
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		))
		msgMap := byteToMap(d.Body)
		msgInfo := byteToMap([]byte(fmt.Sprintf("%v", msgMap["oslo.message"])))
		global.LOG.Debug(fmt.Sprintf("[req-%s]%s msgInfo: %v", reqID, consumer.Tag, msgInfo))
		eventType, exist := msgInfo["event_type"]
		global.LOG.Debug(fmt.Sprintf("[req-%s]%s event_type: %v", reqID, consumer.Tag, eventType))
		if exist {
			if resource, exist := consts.EventStates[eventType.(string)]; exist {
				key := GetResourceMap(resource)
				if publisher, exist := publisherMap[key]; exist {
					parser := NewParser(eventType.(string), resource, key, msgInfo, consumer.IsSDN)
					publishMsg := parser.Parse()
					if publishMsg != nil {
						publishMsg.Timestamp = fmt.Sprintf("%d", time.Now().UnixMilli())
						publishMsg.AzCode = consumer.AzCode
						global.LOG.Info(fmt.Sprintf("[req-%s]%s publishing body (%#v) for event_type %v", reqID, consumer.Tag, publishMsg, eventType))
						if len(buffer) == cap(buffer) {
							global.LOG.Info(fmt.Sprintf("[req-%s]Buffer is full, stop consume rabbitmq message...", reqID))
							buffer <- publishMsg
						}
						go publisher.Call(utils.WrapRequestID(reqID), publishMsg, buffer)
					} else {
						global.LOG.Info(fmt.Sprintf("[req-%s]%s publishing body nil for event_type %v", reqID, consumer.Tag, eventType))
					}
				} else{
					global.LOG.Debug(fmt.Sprintf("[req-%s]Publisher %s not init, don't to send msg.", reqID, key))
				}
			}
		}
		d.Ack(false)
	}
	global.LOG.Info("handle: deliveries channel closed")
}

func byteToMap(body []byte) map[string]interface{} {
	msg := make(map[string]interface{})
	err := json.Unmarshal(body, &msg)
	if err != nil {
		global.LOG.Error("msg unmarshal err", zap.Error(err))
	}
	return msg
}

func GetResourceMap(resource string) (key string) {
	key = resource
	if value, exist := consts.SecondaryMap[resource]; exist {
		key = value
	}
	return
}
