package payment

import (
	"context"
	"net/http"
	"reservation/internal/service/payment"
	"reservation/internal/storage/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type PaymentHandler interface {
	ProcessPayment(c *gin.Context)
	GetPublishableKey(c *gin.Context)
	WebHook(c *gin.Context)
}
type paymentHandler struct {
	logger         *slog.Logger
	srv            payment.PaymentProcessor
	publishableKey string
}

func NewPaymentHandler(logger *slog.Logger, service payment.PaymentProcessor, pkey string) PaymentHandler {
	return &paymentHandler{
		logger:         logger,
		srv:            service,
		publishableKey: pkey,
	}
}

func (p *paymentHandler) ProcessPayment(c *gin.Context) {
	var agent = c.Query("agent")
	var payload db.Reservation

	if err := c.ShouldBindJSON(&payload); err != nil {
		p.logger.Info("unable to bind payment request bosy", err)
		c.JSON(http.StatusOK, err.Error())
		return
	}
	paymentURL, err := p.srv.ProcessPayment(context.Background(), agent, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, paymentURL)

}
func (p *paymentHandler) GetPublishableKey(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"publishableKey": p.publishableKey})
}
func (p *paymentHandler) WebHook(c *gin.Context) {
	p.srv.HandleWebHook(c)
}
