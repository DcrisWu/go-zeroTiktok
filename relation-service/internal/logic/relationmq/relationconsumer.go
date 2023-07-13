package relationmq

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/relation-service/pb/relation"
	"go-zeroTiktok/utils"
	"gorm.io/gorm"
	"k8s.io/klog/v2"
)

var DB *gorm.DB

func RelationConsumer(relationMq *utils.RabbitMq, db *gorm.DB) {
	DB = db
	// 1.声明队列
	_, err := relationMq.Channel.QueueDeclare(
		relationMq.QueueName,
		// 是否持久化
		true,
		// 是否为自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞
		false,
		// 额外属性
		nil,
	)
	if err != nil {
		klog.Info("relation模块声明队列失败")
		panic(err)
	}

	// 2.接收消息
	msgChannel, err := relationMq.Channel.Consume(
		relationMq.QueueName,
		"",
		// 是否自动应答
		true,
		// 是否具有排他性
		false,
		// 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		// 消息队列是否阻塞
		false,
		nil,
	)
	if err != nil {
		klog.Info("relation模块接收消息失败")
	}
	// 3.处理消息
	for msg := range msgChannel {
		var req *relation.ActionReq
		err := json.Unmarshal(msg.Body, &req)
		if err != nil {
			logx.Error("relationMq序列化消费信息失败")
		} else {
			go RelationAction(req)
		}
	}
}

func RelationAction(req *relation.ActionReq) {
	if req.ActionType == 1 {
		err := db.CreateRelation(context.Background(), DB, req.UserId, req.ToUserId)
		if err != nil {
			logx.Error("relationMq点赞失败")
			return
		}
	}
	if req.ActionType == 2 {
		err := db.CancelRelation(context.Background(), DB, req.UserId, req.ToUserId)
		if err != nil {
			logx.Error("relationMq取消点赞失败")
			return
		}
	}
}
