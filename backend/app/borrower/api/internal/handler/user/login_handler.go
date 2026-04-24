package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"zeebee-backend/app/borrower/api/internal/logic/user"
	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/validator"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/response"
)

// login
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				logx.WithContext(r.Context()).
					WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
					Error(err)
			}
		}()
		var req types.UserLoginRequest
		if err = httpx.Parse(r, &req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		v, err := validator.GetEnValidator()
		if err != nil {
			response.InternalServerError(w, errorx.GetErrorMsg(err))
			return
		}
		if err := v.Validate(&req); err != nil {
			response.BadRequest(w, errorx.GetErrorMsg(err))
			return
		}

		l := user.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			// user is unverified
			if e, ok := status.FromError(err); ok && e.Code() == codes.PermissionDenied {
				response.Error(w, 401001, e.Message())
			} else {
				response.BadRequest(w, errorx.GetErrorMsg(err))
			}
			return
		} else {
			response.Success(w, resp)
		}
	}
}
