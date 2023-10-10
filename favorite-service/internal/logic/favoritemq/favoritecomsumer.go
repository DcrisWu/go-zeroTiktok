package favoritemq

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zeroTiktok/favorite-service/pb/favorite"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/utils"
	"gorm.io/gorm"
	"k8s.io/klog/v2"
)

var DB *gorm.DB

func FavoriteConsumer(favoriteMq *utils.RabbitMq, db *gorm.DB) {
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
		// 获取到的消息是amqp.Delivery对象，从中可以获取消息信息
		var req *favorite.FavoriteActionReq
		err := json.Unmarshal(msg.Body, &req)
		if err != nil {
			logx.Error("favoriteMq序列化消费信息失败")
		} else {
			// 重试3次
			err := Retry(func() error {
				return FavoriteAction(req)
			}, 3)
			if err != nil {
				logx.Error("favoriteMq消费失败")
			}
		}
	}
}

func Retry(f func() error, times int) (err error) {
	// 重试机制
	for i := 0; i < times; i++ {
		err = f()
		if err == nil {
			break
		}
	}
	return
}

func FavoriteAction(req *favorite.FavoriteActionReq) error {

	if req.ActionType == 1 {
		// 点赞
		err := db.CreateFavorite(context.Background(), DB, req.UserId, req.VideoId)
		if err != nil {
			return errors.New("favoriteMq添加点赞关系失败")
		}
	}
	if req.ActionType == 2 {
		// 取消点赞
		err := db.CancelFavorite(context.Background(), DB, req.UserId, req.VideoId)
		if err != nil {
			return errors.New("favoriteMq添加点赞关系失败")
		}
	}
	return nil
}
