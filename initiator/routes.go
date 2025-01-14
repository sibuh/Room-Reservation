package initiator

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//register room routes

var roomRoutes = []route{
	{
		path:    "/reserve",
		method:  http.MethodPost,
		handler: roomHandler.Reserve,
	},
	{
		path:    "/add_room",
		method:  http.MethodPost,
		handler: roomHandler.AddRoom,
	},
	{
		path:    "/update_room",
		method:  http.MethodPost,
		handler: roomHandler.UpdateRoom,
	},
	{
		path:    "/reservations",
		method:  http.MethodPost,
		handler: roomHandler.GetRoomReservations,
	},
}

//register all user related routes

var userRoutes = []route{
	{
		path:    "/signup",
		method:  http.MethodPost,
		handler: userHandler.Signup,
	},
	{
		path:        "/login",
		method:      http.MethodPost,
		handler:     userHandler.Login,
		middlewares: []gin.HandlerFunc{},
	},
	{
		path:    "/refresh",
		method:  http.MethodGet,
		handler: userHandler.Refresh,
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
		handler: hotelHandler.Register,
		middlewares: []gin.HandlerFunc{
			mw.Authorize(),
		},
	},
	{
		path:    "/search",
		method:  http.MethodPost,
		handler: hotelHandler.SearchHotel,
	},
	{
		path:    "/hotel",
		method:  http.MethodGet,
		handler: hotelHandler.GetHotelByName,
	},
	{
		path:    "/hotels",
		method:  http.MethodGet,
		handler: hotelHandler.GetHotels,
	},
}

var paymentRoutes = []route{
	{
		path:    "/pkey",
		method:  http.MethodGet,
		handler: paymentHandler.GetPublishableKey,
	},
	{
		path:    "/webhook",
		method:  http.MethodPost,
		handler: paymentHandler.WebHook,
	},
	{
		path:    "/pay",
		method:  http.MethodPost,
		handler: paymentHandler.ProcessPayment,
	},
}
