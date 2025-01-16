package initiator

import (
	"net/http"
	"reservation/internal/handler/hotel"
	"reservation/internal/handler/middleware"
	"reservation/internal/handler/payment"
	"reservation/internal/handler/room"
	"reservation/internal/handler/user"

	"github.com/gin-gonic/gin"
)

func ListRoutes(
	rh room.RoomHandler,
	hh hotel.HotelHandler,
	uh user.UserHandler,
	ph payment.PaymentHandler,
	mw middleware.Middleware) [][]route {
	//register room routes

	var roomRoutes = []route{
		{
			path:    "/reserve",
			method:  http.MethodPost,
			handler: rh.Reserve,
		},
		{
			path:    "/add_room",
			method:  http.MethodPost,
			handler: rh.AddRoom,
		},
		{
			path:    "/update_room",
			method:  http.MethodPost,
			handler: rh.UpdateRoom,
		},
		{
			path:    "/reservations",
			method:  http.MethodPost,
			handler: rh.GetRoomReservations,
		},
	}

	//register all user related routes

	var userRoutes = []route{
		{
			path:    "/signup",
			method:  http.MethodPost,
			handler: uh.Signup,
		},
		{
			path:        "/login",
			method:      http.MethodPost,
			handler:     uh.Login,
			middlewares: []gin.HandlerFunc{},
		},
		{
			path:    "/refresh",
			method:  http.MethodGet,
			handler: uh.Refresh,
			middlewares: []gin.HandlerFunc{
				mw.Authorize(),
			},
		},
	}

	//register all hotel related routes

	var hotelRoutes = []route{
		{
			path:    "/register",
			method:  http.MethodPost,
			handler: hh.Register,
			middlewares: []gin.HandlerFunc{
				mw.Authorize(),
			},
		},
		{
			path:    "/search",
			method:  http.MethodPost,
			handler: hh.SearchHotel,
		},
		{
			path:    "/hotel",
			method:  http.MethodGet,
			handler: hh.GetHotelByName,
		},
		{
			path:    "/hotels",
			method:  http.MethodGet,
			handler: hh.GetHotels,
		},
	}

	var paymentRoutes = []route{
		{
			path:    "/pkey",
			method:  http.MethodGet,
			handler: ph.GetPublishableKey,
		},
		{
			path:    "/webhook",
			method:  http.MethodPost,
			handler: ph.WebHook,
		},
		{
			path:    "/pay",
			method:  http.MethodPost,
			handler: ph.ProcessPayment,
		},
	}
	allRoutes := [][]route{
		userRoutes,
		hotelRoutes,
		roomRoutes,
		paymentRoutes,
	}

	return allRoutes
}
