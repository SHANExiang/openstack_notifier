package consume

import (
	"fmt"
	"github.com/streadway/amqp"
	"sincerecloud.com/openstack_notifier/consts"
	"sincerecloud.com/openstack_notifier/global"
	"sincerecloud.com/openstack_notifier/services/publish"
	"sync"
)

type RabbitMQ struct {
	Conn    	*amqp.Connection
	Channel  	*amqp.Channel
	QueueName 	string
	Exchange 	string
	Key 		string
	AmqpURIPool []string
	wg          sync.WaitGroup
	AzCode      string
	Tag         string
	mu          sync.Mutex
	Buffers     []chan *publish.PublishMsg
	Running     bool
	IsSDN       bool
}

func NewRabbitMQ(queueName, exchange, key string, amqpURIPool []string) (obj *RabbitMQ) {
	obj = &RabbitMQ{
		QueueName: queueName,
		Exchange: exchange,
		Key: key,
		AmqpURIPool: amqpURIPool,
	}
	return
}

func (mq *RabbitMQ) GetAmqpURI() string {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	// Round Robin
	uri := mq.AmqpURIPool[0]
	mq.AmqpURIPool = append(mq.AmqpURIPool[1:], uri)
	return uri
}

func (mq *RabbitMQ) GetChannel() error {
	var err error
	var conn *amqp.Connection
	for i := 0;i < len(mq.AmqpURIPool);i++ {
		uri := mq.GetAmqpURI()
		global.LOG.Info(fmt.Sprintf("consumer dialing %q", uri))
		conn, err = amqp.Dial(uri)
		if err != nil {
			global.LOG.Error(fmt.Sprintf("%s dial: %s", uri, err))
		} else {
			break
		}
	}
	if conn == nil && err != nil {
		return err
	}

	mq.Conn = conn
	global.LOG.Info(fmt.Sprintf("%s got Connection, getting Channel", mq.Tag))
	channel, err := conn.Channel()
	if err != nil {
		global.LOG.Error(fmt.Sprintf("%s channel err: %s", mq.Tag, err))
		return err
	}
	mq.Channel = channel
	return nil
}

func (mq *RabbitMQ) ExchangeDeclare() error {
	if err := mq.Channel.ExchangeDeclare(
		mq.Exchange,           // name
		consts.ExchangeType,   // type
		true,          // durable
		false,      // auto-deleted
		false,        // internal
		false,        // noWait
		nil,            // arguments
	); err != nil {
		global.LOG.Error(fmt.Sprintf("%s exchange declare err: %s", mq.Tag, err))
		return err
	}
	return nil
}

func (mq *RabbitMQ) QueueDeclare(queueName string) (amqp.Queue, error) {
	queue, err := mq.Channel.QueueDeclare(
		queueName, 			    // name of the queue
		false,           // durable
		false,    	// delete when unused
		false,    	    // exclusive
		false,    	    // noWait
		nil,       		// arguments
	)
	return queue, err
}

func (mq *RabbitMQ) Close() {
	mq.Running = false
	if mq.Channel != nil {
		global.LOG.Info(fmt.Sprintf("%s close channel", mq.Tag))
		err := mq.Channel.Close()
		if err != nil {
			global.LOG.Error(fmt.Sprintf("%s close channel err:%s", mq.Tag, err))
		}
	}
	if mq.Conn != nil {
		global.LOG.Info(fmt.Sprintf("%s close connection", mq.Tag))
		err := mq.Conn.Close()
		if err != nil {
			global.LOG.Error(fmt.Sprintf("%s close channel err:%s", mq.Tag, err))
		}
	}
}
