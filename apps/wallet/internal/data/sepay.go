package data

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"vn.vato.zora.be.api/apps/wallet/internal/biz"
	"vn.vato.zora.be.api/apps/wallet/internal/conf"
)

type sepayGateway struct {
	apiKey        string
	webhookSecret string
	log           *log.Helper
}

func NewSepayGateway(c *conf.Data, logger log.Logger) biz.PaymentGateway {
	var apiKey, secret string
	if c != nil && c.Sepay != nil {
		apiKey = c.Sepay.ApiKey
		secret = c.Sepay.WebhookSecret
	}
	return &sepayGateway{
		apiKey:        apiKey,
		webhookSecret: secret,
		log:           log.NewHelper(logger),
	}
}

func (g *sepayGateway) CreateQR(ctx context.Context, orderID string, amount int64) (string, time.Time, error) {
	// TODO: POST https://api.sepay.vn/v1/qr/create
	//   header: Authorization: Bearer <apiKey>
	g.log.Infof("Sepay.CreateQR: orderID=%s amount=%d", orderID, amount)
	qrURL := fmt.Sprintf("https://qr.sepay.vn/img?bank=MB&acc=0123456789&template=compact&amount=%d&des=%s", amount, orderID)
	return qrURL, time.Now().Add(15 * time.Minute), nil
}

// VerifyWebhookSignature xác thực chữ ký HMAC-SHA256 Sepay gửi về.
// Sepay ký: HMAC-SHA256(webhookSecret, orderID + "|" + amount + "|" + status)
func (g *sepayGateway) VerifyWebhookSignature(_ context.Context, payload *biz.SepayWebhookPayload) error {
	if g.webhookSecret == "" {
		return fmt.Errorf("webhook secret not configured")
	}
	message := fmt.Sprintf("%s|%d|%s", payload.OrderID, payload.Amount, payload.Status)
	mac := hmac.New(sha256.New, []byte(g.webhookSecret))
	mac.Write([]byte(message))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(payload.Signature)) {
		return fmt.Errorf("signature mismatch")
	}
	return nil
}
