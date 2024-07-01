package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lifememo/application/message/api/internal/logic"
	"lifememo/application/message/api/internal/svc"
	"lifememo/application/message/api/internal/types"
)

func DeleteMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteMessageRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDeleteMessageLogic(r.Context(), svcCtx)
		resp, err := l.DeleteMessage(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
