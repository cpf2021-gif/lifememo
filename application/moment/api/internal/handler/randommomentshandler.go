package handler

import (
	"net/http"

	"lifememo/application/moment/api/internal/logic"
	"lifememo/application/moment/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RandomMomentsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewRandomMomentsLogic(r.Context(), svcCtx)
		resp, err := l.RandomMoments()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
