package payments

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	InvalidSignature = errors.New("Signature is Invalid")
)

type PaymentService interface {
	HandleWebhookEvent(ctx context.Context, rawbody []byte, signatureHeader string) error
}
type paymentService struct {
	repo          PaymentRepository
	redis         *redis.Client
	webhookSecret string
}

func NewPaymentService(
	repo PaymentRepository,
	redis *redis.Client,
	webhookSecret string,
) PaymentService {
	return &paymentService{
		repo:          repo,
		redis:         redis,
		webhookSecret: webhookSecret,
	}
}
func (s *paymentService) HandleWebhookEvent(ctx context.Context, rawbody []byte, signatureHeader string) error {
	var event RazorpayWebhookEvent
	receivedSignatureBytes, err := hex.DecodeString(signatureHeader)
	if err != nil {
		return err
	}

	mac := hmac.New(sha256.New, []byte(s.webhookSecret))
	mac.Write(rawbody)
	expectedSignature := mac.Sum(nil)

	val := hmac.Equal(expectedSignature, receivedSignatureBytes)
	if !val {
		return InvalidSignature
	}
	err = json.Unmarshal(rawbody, &event)
	if err != nil {
		return err
	}
	if event.Event != "payment.captured" {
		return nil
	} else {
		bookingIDStr := event.Payload.Payment.Entity.Notes.BookingID
		eventID := event.Payload.Payment.Entity.Notes.EventID
		seatNumber := event.Payload.Payment.Entity.Notes.SeatNumber
		userID := event.Payload.Payment.Entity.Notes.UserID
		eventId, err := strconv.ParseInt(eventID, 10, 64)
		if err != nil {
			return nil
		}

		user_id, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			return nil
		}
		bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
		if err != nil {
			return err
		}
		err = s.repo.ConfirmPayment(ctx, bookingID, event.Payload.Payment.Entity.ID, event.Payload.Payment.Entity.OrderID)
		if err != nil {
			return err
		}
		seatHoldKey := fmt.Sprintf("hold:%d:%s", eventId, seatNumber)
		userHoldKey := fmt.Sprintf("user:hold:%d", user_id)
		s.redis.Del(ctx, seatHoldKey, userHoldKey)
	}
	return nil
}
