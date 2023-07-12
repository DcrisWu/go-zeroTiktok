package mq

import (
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/logx"
)

type RabbitMq struct {
	conn       *amqp.Connection
	Channel    *amqp.Channel
	QueueName  string
	Exchange   string
	RoutingKey string
	MqUrl      string
}

var Rmq *RabbitMq

// InitRabbitMQ 初始化RabbitMQ的连接和通道。
func InitRabbitMQ(mqUrl string) {
	Rmq = &RabbitMq{
		MqUrl: mqUrl,
	}
	dial, err := amqp.Dial(mqUrl)
	if err != nil {
		logx.Error(err)
		return
	}
	Rmq.conn = dial
}

func NewRabbitMq(queueName, exchange, routingKey string) *RabbitMq {
	rabbitMQ := RabbitMq{
		QueueName:  queueName,
		conn:       Rmq.conn,
		Exchange:   exchange,
		RoutingKey: routingKey,
		MqUrl:      Rmq.MqUrl,
	}
	var err error
	// 创建Channel
	rabbitMQ.Channel, err = rabbitMQ.conn.Channel()
	if err != nil {
		logx.Error(err)
	}
	return &rabbitMQ
}

// ReleaseRes 关闭mq通道和mq的连接。
func (r *RabbitMq) ReleaseRes() {
	err := r.conn.Close()
	if err != nil {
		logx.Error(err)
		return
	}
}
