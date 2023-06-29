package feed

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zeroTiktok/system-api/internal/logic/feed"
	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"
)

func GetVideoListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FeedReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := feed.NewGetVideoListLogic(r.Context(), svcCtx)
		resp, err := l.GetVideoList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
