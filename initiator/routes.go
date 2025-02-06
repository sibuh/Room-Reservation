package initiator

import (
	"net/http"
	"reservation/internal/handler/hotel"
	"reservation/internal/handler/middleware"
	"reservation/internal/handler/payment"
	"reservation/internal/handler/room"
	roomtype "reservation/internal/handler/room_type"
	"reservation/internal/handler/user"

	"github.com/gin-gonic/gin"
)

func ListRoutes(
	rh room.RoomHandler,
	hh hotel.HotelHandler,
	uh user.UserHandler,
	ph payment.PaymentHandler,
	rth roomtype.RoomTypeHandler,
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
			method:  http.MethodGet,
			handler: rh.GetRoomReservations,
		},
		{
			path:    "/:hotel_id/rooms",
			method:  http.MethodGet,
			handler: rh.GetHotelRooms,
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
		{
			path:    "/hotel/:hotel_id",
			method:  http.MethodPatch,
			handler: hh.VerifyHotel,
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
	var roomTypeRoutes = []route{
		{
			path:    "/room_type",
			method:  http.MethodPost,
			handler: rth.CreateRoomType,
			//TODO: middleware
			//this endpoint has to be accessed only by super admin
			//to add all possible types of room types
		},
		{
			path:    "/room_types",
			method:  http.MethodGet,
			handler: rth.GetRoomTypes,
			//TODO: middleware
			//this endpoint has to be accessed only by super admin
			//to add all possible types of room types
		},
	}
	allRoutes := [][]route{
		userRoutes,
		hotelRoutes,
		roomRoutes,
		paymentRoutes,
		roomTypeRoutes,
	}

	return allRoutes
}

func RegisterRoutes(g *gin.RouterGroup, routes []route) {
	for _, route := range routes {
		route.middlewares = append(route.middlewares, route.handler)
		g.Handle(route.method, route.path, route.middlewares...)
	}
}
