package handler

import (
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/response"
)

// stripeWebhookMaxBodyBytes caps the webhook payload size; Stripe events are
// small, anything larger is rejected before buffering.
const stripeWebhookMaxBodyBytes = 64 * 1024

// RegisterStripeWebhookRoute registers the public Stripe webhook endpoint as a
// hand-written extra route: signature verification needs the raw request body,
// which goctl-generated handlers cannot expose, and extra routes bypass the
// AuthMiddleware (Stripe authenticates via the signature header instead).
func RegisterStripeWebhookRoute(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/webhook/stripe",
		Handler: stripeWebhookHandler(serverCtx),
	})
}

func stripeWebhookHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := io.ReadAll(http.MaxBytesReader(w, r.Body, stripeWebhookMaxBodyBytes))
		if err != nil {
			response.BadRequest(w, "invalid webhook payload")
			return
		}

		_, err = serverCtx.StripeWebhookRpcClient.Handle(r.Context(), &admin.StripeWebhookRequest{
			Payload:   payload,
			Signature: r.Header.Get("Stripe-Signature"),
		})
		if err != nil {
			logx.WithContext(r.Context()).WithFields(
				logx.Field("path", r.URL.Path),
			).Errorf("stripe webhook failed: %v", err)
			// 400 stops Stripe from retrying (bad signature); other errors
			// answer 500 so Stripe retries with backoff
			if status.Code(err) == codes.InvalidArgument {
				response.BadRequest(w, "invalid stripe signature")
				return
			}
			response.InternalServerError(w, "webhook processing failed")
			return
		}
		response.Success(w, nil)
	}
}
