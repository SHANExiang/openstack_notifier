package services

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"reflect"
	"sincerecloud.com/openstack_notifier/nonick"
	"sincerecloud.com/openstack_notifier/services/consume"
	"testing"
)

func TestDial(t *testing.T) {
	url := "amqp://openstack:a0va51Upnwwm4nE0Vg9Ba5f6J1kwoMexLkHnrIbh@10.50.1.204:5672/"
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("%s dial: %s\n", url, err)
	}
	if conn == nil {
		log.Fatalln("Failed to established connection")
	}
}

func TestInitInitConsumer(t *testing.T) {
	mqInfo := nonick.MQInfo{
		AzCode: "cn-Guizhou-1-D",
		Endpoint: "http://10.50.1.200",
		MqChannel: "ecs_voneyun_topic",
		MqEndpoint: "10.50.1.1:5672,10.50.1.2:5672, 10.50.1.3:5672",
		MqPassword: "a0NvYRhEdtvYtuIww8k5QX145l8mV6vszmu4VPuq",
		MqUserName: "openstack",
		Name: "RegionOne",
		NickName: "RegionF"}
	amqpURIPool := constructAMQPURI(mqInfo)
	fmt.Println(amqpURIPool)
	consumer := consume.NewRabbitMQ(
		mqInfo.MqChannel,
		mqInfo.MqChannel,
		"#",
		amqpURIPool)
	if consumer.GetAmqpURI() != "amqp://openstack:a0NvYRhEdtvYtuIww8k5QX145l8mV6vszmu4VPuq@10.50.1.1:5672/" {
		log.Fatalln("Failed to get amqp uri 1")
	}
	if consumer.GetAmqpURI() != "amqp://openstack:a0NvYRhEdtvYtuIww8k5QX145l8mV6vszmu4VPuq@10.50.1.2:5672/" {
		log.Fatalln("Failed to get amqp uri 2")
	}
	if consumer.GetAmqpURI() != "amqp://openstack:a0NvYRhEdtvYtuIww8k5QX145l8mV6vszmu4VPuq@10.50.1.3:5672/" {
		log.Fatalln("Failed to get amqp uri 3")
	}
	if consumer.GetAmqpURI() != "amqp://openstack:a0NvYRhEdtvYtuIww8k5QX145l8mV6vszmu4VPuq@10.50.1.1:5672/" {
		log.Fatalln("Failed to get amqp uri 1")
	}
}

func TestName(t *testing.T) {
	content := []byte(`{"status":100,"code":"1-100","msg":"操作成功","data":[{"mq_channel":null,"mq_user_name":"openstack","mq_password":"a05GgaHLpstzvEii2jQnI251yqCDcMW0mMeSn0wD","mq_endpoint":"10.50.114.157:5672","nick_name":"RegionB","name":"RegionOne","endpoint":"http://10.50.114.157","az_code":"cn-Guizhou-2-B","sdn":false},{"mq_channel":null,"mq_user_name":null,"mq_password":null,"mq_endpoint":null,"nick_name":"RegionC","name":"RegionOne","endpoint":"http://10.50.114.112","az_code":"cn-Guizhou-1-C","sdn":false},{"mq_channel":null,"mq_user_name":null,"mq_password":null,"mq_endpoint":null,"nick_name":"Region-2","name":"Region-2","endpoint":"http://10.50.114.123","az_code":"cn-Guizhou-1-B","sdn":false},{"mq_channel":"","mq_user_name":"","mq_password":"","mq_endpoint":"","nick_name":"RegionE","name":"RegionOne","endpoint":"http://10.50.1.57","az_code":"cn-Guizhou-2-A","sdn":false},{"mq_channel":"ecs_voneyun_topic","mq_user_name":"openstack","mq_password":"a0NvYRhEdtvYtuIww8k5QX145l8mV6vszmu4VPuq","mq_endpoint":"10.50.1.1:5672,10.50.1.2:5672,10.50.1.3:5672","nick_name":"RegionF","name":"RegionOne","endpoint":"http://10.50.1.200","az_code":"cn-Guizhou-1-D","sdn":false},{"mq_channel":"ecs_voneyun_topic","mq_user_name":"openstack","mq_password":"a0hOPSbf4AydvrRgtZMiHLWug6oDfFQGH4U8WuSp","mq_endpoint":"10.50.31.1:5672","nick_name":"RegionOne","name":"RegionOne","endpoint":"http://10.50.31.1","az_code":"cn-Guizhou-1-E","sdn":true},{"mq_channel":"ecs_voneyun_topic","mq_user_name":"openstack","mq_password":"a01qmGDuC5y79sm0QoWLw44zbnB65DUer3JGw9rm","mq_endpoint":"10.251.28.41:5672","nick_name":"RegionG","name":"RegionOne","endpoint":"http://10.251.28.41","az_code":"cn-Guizhou-1-F","sdn":true}],"resultType":1,"elapsedTime":0,"timestamp":0,"exception":null,"traceId":null,"requestId":null}`)
	var resp nonick.MqInfoResp
	err := json.Unmarshal(content, &resp)
	if err != nil || reflect.DeepEqual([]nonick.MQInfo{}, resp.Data) {
		log.Fatalln("Failed to parse json content", err)
	}
}
