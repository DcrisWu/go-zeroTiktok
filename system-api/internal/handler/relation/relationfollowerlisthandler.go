package relation

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zeroTiktok/system-api/internal/logic/relation"
	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"
)

func RelationFollowerListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RelationFollowerListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := relation.NewRelationFollowerListLogic(r.Context(), svcCtx)
		resp, err := l.RelationFollowerList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
