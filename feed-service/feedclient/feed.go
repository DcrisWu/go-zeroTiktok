// Code generated by goctl. DO NOT EDIT.
// Source: feed.proto

package feedclient

import (
	"context"

	"go-zeroTiktok/feed-service/pb/feed"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	FeedReq  = feed.FeedReq
	FeedResp = feed.FeedResp
	User     = feed.User
	Video    = feed.Video

	Feed interface {
		Feed(ctx context.Context, in *FeedReq, opts ...grpc.CallOption) (*FeedResp, error)
	}

	defaultFeed struct {
		cli zrpc.Client
	}
)

func NewFeed(cli zrpc.Client) Feed {
	return &defaultFeed{
		cli: cli,
	}
}

func (m *defaultFeed) Feed(ctx context.Context, in *FeedReq, opts ...grpc.CallOption) (*FeedResp, error) {
	client := feed.NewFeedClient(m.cli.Conn())
	return client.Feed(ctx, in, opts...)
}
