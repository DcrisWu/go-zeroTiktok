package mq

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zeroTiktok/favorite-service/pb/favorite"
	"go-zeroTiktok/models/db"
	"gorm.io/gorm"
	"k8s.io/klog/v2"
)

var DB *gorm.DB

func FavoriteConsumer(favoriteMq *RabbitMq, db *gorm.DB) {
	DB = db
	_, err := favoriteMq.Channel.QueueDeclare(favoriteMq.QueueName, true, false, false, false, nil)
	if err != nil {
		klog.Fatalf("favorite add consumer declare error")
		panic(err)
	}

	//2、接收消息
	msgChanel, err := favoriteMq.Channel.Consume(
		favoriteMq.QueueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//消息队列是否阻塞
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	for msg := range msgChanel {
		// 这里写你的处理逻辑
		// 获取到的消息是amqp.Delivery对象，从中可以获取消息信息
		go FavoriteAction(string(msg.Body))
	}
}

func FavoriteAction(msg string) {
	var req *favorite.FavoriteActionReq
	err := json.Unmarshal([]byte(msg), &req)
	if err != nil {
		logx.Error("favoriteMq序列化消费信息失败")
		return
	}
	if req.ActionType == 1 {
		// 点赞
		err := db.CreateFavorite(context.Background(), DB, req.UserId, req.VideoId)
		if err != nil {
			logx.Errorf("favoriteMq添加点赞关系失败")
			return
		}
	}
	if req.ActionType == 2 {
		// 取消点赞
		err := db.CancelFavorite(context.Background(), DB, req.UserId, req.VideoId)
		if err != nil {
			logx.Errorf("favoriteMq取消点赞关系失败")
			return
		}
	}
}
