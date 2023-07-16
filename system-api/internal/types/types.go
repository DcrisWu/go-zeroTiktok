// Code generated by goctl. DO NOT EDIT.
package types

type RegisterReq struct {
	UserName string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

type RegisterResp struct {
	StatusCode int64  `json:"status_code"` // 状态码 0-成功 1-失败
	StatusMsg  string `json:"status_msg"`  // 状态信息
	UserId     int64  `json:"user_id"`     // 用户id
	Token      string `json:"token"`       // token
}

type LoginReq struct {
	UserName string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

type LoginResp struct {
	StatusCode int64  `json:"status_code"` // 状态码 0-成功 1-失败
	StatusMsg  string `json:"status_msg"`  // 状态信息
	UserId     int64  `json:"user_id"`     // 用户id
	Token      string `json:"token"`       // token
}

type User struct {
	Id            int64  `json:"id"`             // 用户id
	Name          string `json:"name"`           // 用户名
	FollowCount   int64  `json:"follow_count"`   // 关注数
	FollowerCount int64  `json:"follower_count"` // 粉丝数
	IsFollow      bool   `json:"is_follow"`      // 是否关注
}

type UserReq struct {
	UserId int64 `form:"user_id"`
}

type UserResp struct {
	StatusCode int64  `json:"status_code"` // 状态码 0-成功 1-失败
	StatusMsg  string `json:"status_msg"`  // 状态信息
	User       *User  `json:"user"`        // 用户信息
}

type FeedReq struct {
	Authorization string `header:"Authorization,optional"` // token
	LatestTime    string `json:"latest_time,optional"`     // 最新视频的时间，格式 yyyy-mm-dd hh:mm:ss
}

type Video struct {
	Id            int64  `json:"id"`
	Author        *User  `json:"author"`
	PlayUrl       string `json:"play_url"`       // 播放地址
	CoverUrl      string `json:"cover_url"`      // 封面地址
	FavoriteCount int64  `json:"favorite_count"` // 收藏数
	CommentCount  int64  `json:"comment_count"`  // 评论数
	IsFavorite    bool   `json:"is_favorite"`    // 是否收藏
	Title         string `json:"title"`          // 视频标题
}

type FeedResp struct {
	StatsCode int64    `json:"stats_code"`
	StatusMsg string   `json:"status_msg"`
	NextTime  int64    `json:"next_time"`
	VideoList []*Video `json:"video_list"`
}

type DataInfo struct {
	PlayUrl  string `json:"play_url"`
	CoverUrl string `json:"cover_url"`
}

type PublishActionReq struct {
	Data  DataInfo `json:"data"`
	Title string   `json:"title"`
}

type PublishActionResp struct {
	StatusCode int64  `json:"status_code"` // 状态码 1-成功 2-失败
	StatusMsg  string `json:"status_msg"`
}

type PublishListReq struct {
	UserId int64 `form:"user_id"`
}

type PublishListResp struct {
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
	VideoList  []*Video `json:"video_list"`
}

type Comment struct {
	Id         int64  `json:"id"`
	UserInfo   *User  `json:"user"`        // 评论用户信息
	Content    string `json:"content"`     // 评论内容
	CreateDate string `json:"create_date"` // 评论发布日期，格式 mm-dd
}

type CommentActionReq struct {
	VideoId     string `json:"vedio_id"`              // 视频id
	ActionType  string `json:"action_type"`           // 1-发布评论，2-删除评论
	CommentText string `json:"comment_text,optional"` // 可选，用户填写的评论内容，在action_type=1的时候使用
	CommentId   string `json:"comment_id,optional"`   // 可选，要删除的评论id，在action_type=2的时候使用
}

type CommentActionResp struct {
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
	CommentObj *Comment `json:"comment,optional"` // 多选一且必须只能符合下列其中一组子节点（即XOR，
}

type CommentListReq struct {
	VedioId string `form:"vedioId"`
}

type CommentListResp struct {
	StatusCode  int64      `json:"status_code"`
	StatusMsg   string     `json:"status_msg"`
	CommentList []*Comment `json:"array"`
}

type FavoriteActionReq struct {
	VideoId    string `json:"video_id"`
	ActionType string `json:"action_type"` // 1-点赞，2-取消点赞
}

type FavoriteActionResp struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`
}

type FavoriteListReq struct {
	UserId string `path:"user_id"`
}

type FavoriteListResp struct {
	StatusCode int64    `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string   `json:"status_msg"`  // 返回状态描述
	VideoList  []*Video `json:"video_list"`  // 户点赞视频列表
}

type RelationActionReq struct {
	ToUserId   string `json:"to_user_id"`  // 被关注用户id
	ActionType string `json:"action_type"` // 1-关注，2-取消关注
}

type RelationActionResp struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

type RelationFollowListReq struct {
	UserId string `path:"user_id"`
}

type RelationFollowListResp struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg"`  // 返回状态描述
	UserList   []*User `json:"user_list"`   // 用户关注列表
}

type RelationFollowerListReq struct {
	UserId string `path:"user_id"`
}

type RelationFollowerListResp struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg"`  // 返回状态描述
	UserList   []*User `json:"user_list"`   // 用户粉丝列表
}
