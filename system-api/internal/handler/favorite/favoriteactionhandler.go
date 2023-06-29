package favorite

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zeroTiktok/system-api/internal/logic/favorite"
	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"
)

func FavoriteActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FavoriteActionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := favorite.NewFavoriteActionLogic(r.Context(), svcCtx)
		resp, err := l.FavoriteAction(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
