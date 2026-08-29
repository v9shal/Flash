package payments

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service PaymentService
}

func NewPaymenthandler(service PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	rawbody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sig := c.GetHeader("X-Razorpay-Signature")
	if sig == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "signature header missing"})
		return
	}
	err = h.service.HandleWebhookEvent(c.Request.Context(), rawbody, sig)
	if errors.Is(err, InvalidSignature) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "Payment recevied"})
}
