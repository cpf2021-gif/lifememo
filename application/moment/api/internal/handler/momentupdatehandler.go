package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lifememo/application/moment/api/internal/logic"
	"lifememo/application/moment/api/internal/svc"
	"lifememo/application/moment/api/internal/types"
)

func MomentUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MomentUpdateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewMomentUpdateLogic(r.Context(), svcCtx)
		resp, err := l.MomentUpdate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
