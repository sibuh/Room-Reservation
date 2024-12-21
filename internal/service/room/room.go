package room

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"reservation/internal/storage/db"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type RoomService interface {
	ReserveRoom(ctx context.Context, param ReserveRoom) (string, error)
	UpdateRoom(ctx context.Context, param UpdateRoom) (Room, error)
}

type roomService struct {
	db.Querier
	url    string
	logger slog.Logger
}

func NewRoomService(q db.Querier, url string) RoomService {
	return &roomService{
		Querier: q,
		url:     url,
	}
}

func (rs *roomService) ReserveRoom(ctx context.Context, param ReserveRoom) (string, error) {
	if err := param.Validate(); err != nil {
		return "", ErrInvalidInput
	}
	_, err := rs.Querier.UpdateRoom(ctx, db.UpdateRoomParams{
		UserID: uuid.NullUUID{
			UUID:  param.UserID,
			Valid: true,
		},

		ID: param.RoomID,
	})
	if err != nil {
		return "", ErrReservationFailed
	}
	req := CheckoutRequest{
		ProductID:   param.RoomID.String(),
		CallbackURL: "http://localhost:9090/callback", //TODO: url should be read from config
	}
	ssn, err := rs.createCheckoutSession(ctx, req)
	if err != nil {
		return "", ErrCheckoutSessionFailed
	}
	return ssn.PaymentURL, nil

}
func (rs *roomService) createCheckoutSession(ctx context.Context, req CheckoutRequest) (CheckoutResponse, error) {
	bbyte, err := json.Marshal(req)
	if err != nil {
		return CheckoutResponse{}, err
	}

	request, err := http.NewRequest(http.MethodPost, rs.url, bytes.NewBuffer(bbyte))
	if err != nil {
		return CheckoutResponse{}, err
	}
	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return CheckoutResponse{}, err
	}
	var session CheckoutResponse
	if err := json.NewDecoder(res.Body).Decode(&session); err != nil {
		return CheckoutResponse{}, err
	}
	return session, nil
}
func (rs *roomService) UpdateRoom(ctx context.Context, param UpdateRoom) (Room, error) {
	rm, err := rs.Querier.UpdateRoom(ctx, db.UpdateRoomParams{
		Status: db.RoomStatus(param.Status),
		UserID: uuid.NullUUID{
			UUID:  param.UserID,
			Valid: true,
		},
		ID: param.ID,
	})
	if err != nil {
		rs.logger.Error("failed to update room", err)
		return Room{}, err
	}
	return Room{
		ID:         rm.ID,
		RoomNumber: rm.RoomNumber,
		UserID:     rm.UserID.UUID,
		HotelID:    rm.HotelID,
		CreatedAt:  rm.CreatedAt,
		UpdatedAt:  rm.UpdatedAt,
	}, nil
}

// import (
// 	"encoding/json"
// 	"bus_ticket/internal/handler"
// 	"bus_ticket/internal/model"
// 	"bus_ticket/internal/module"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stripe/stripe-go/v78"
// 	"golang.org/x/exp/slog"
// )

// type payment struct {
// 	publishableKey string
// 	secretKey      string
// 	logger         *slog.Logger
// 	pm             module.Payment
// }

// func Init(pkey, secretKey string, logger *slog.Logger, pm module.Payment) handler.Payment {
// 	return &payment{
// 		publishableKey: pkey,
// 		secretKey:      secretKey,
// 		logger:         logger,
// 		pm:             pm,
// 	}
// }

// func (p *payment) GetPublishableKey(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"publishableKey": p.publishableKey})
// }

// func (p *payment) HandleCreatePaymentIntent(c *gin.Context) {

// 	stripe.Key = p.secretKey
// 	eventID, _ := strconv.ParseInt(c.Params.ByName("id"), 10, 32)
// 	user := c.Value("user").(model.User)
// 	clientSecret, err := p.pm.CreatePaymentIntent(c, int32(user.ID), int32(eventID))
// 	if err != nil {
// 		newError := err.(*model.Error)
// 		c.JSON(newError.ErrCode, newError)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"clientSecret": clientSecret,
// 	})
// }
// func (p *payment) PaymentWebhook(c *gin.Context) {
// 	var stripeEvent stripe.Event
// 	if err := c.ShouldBindJSON(&stripeEvent); err != nil {
// 		newError := model.Error{
// 			ErrCode:   http.StatusOK,
// 			Message:   "failed to bind request body",
// 			RootError: err,
// 		}
// 		p.logger.Info("unable to bind event request bosy", newError)
// 		c.JSON(newError.ErrCode, newError)
// 		return
// 	}
// 	switch stripeEvent.Type {
// 	case "payment_intent.succeeded":
// 		var paymentIntent stripe.PaymentIntent
// 		err := json.Unmarshal(stripeEvent.Data.Raw, &paymentIntent)
// 		if err != nil {
// 			newError := model.Error{
// 				ErrCode:   http.StatusBadRequest,
// 				Message:   "failed to unmarshal event data to payment intent",
// 				RootError: err,
// 			}
// 			p.logger.Error("Error parsing webhook JSON", newError)
// 			c.JSON(newError.ErrCode, newError)
// 			return
// 		}
// 		p.logger.Info("PaymentIntent was successful!")
// 	case "payment_method.attached":
// 		var paymentMethod stripe.PaymentMethod
// 		err := json.Unmarshal(stripeEvent.Data.Raw, &paymentMethod)
// 		if err != nil {
// 			newError := model.Error{
// 				ErrCode:   http.StatusBadRequest,
// 				Message:   "failed to unmarshal event data to stripe paymentMethod object",
// 				RootError: err,
// 			}
// 			p.logger.Error("Error parsing webhook JSON", newError)
// 			c.JSON(newError.ErrCode, newError)
// 			return
// 		}
// 		p.logger.Info("PaymentMethod was attached to a Customer!")

// 	default:
// 		p.logger.Info("unhandled envet type", stripeEvent.Type)
// 	}

// 	c.JSON(http.StatusOK, nil)

// }
