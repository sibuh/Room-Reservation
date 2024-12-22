package initiator

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"reservation/internal/service/hotel"
	"reservation/internal/service/room"
	"reservation/internal/service/user"

	hh "reservation/internal/handler/hotel"
	"reservation/internal/handler/middleware"
	rh "reservation/internal/handler/room"
	uh "reservation/internal/handler/user"

	"reservation/internal/storage/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"golang.org/x/exp/slog"
)

type route struct {
	path        string
	method      string
	handler     gin.HandlerFunc
	middlewares []gin.HandlerFunc
}

func Initiate() {
	//create connection to database
	conn, err := pgx.Connect(context.Background(), "")
	if err != nil {
		log.Fatal("failed to create connection to db", err)
	}
	//create logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	// initialize storage layer
	queries := db.New(conn)
	key := "1234567890sdfghjkjjkhk"
	duration := 1 * time.Minute
	//initialize services
	userService := user.NewUserService(logger, queries, key, duration)
	roomService := room.NewRoomService(queries, "url")
	hotelService := hotel.NewHotelService(queries, logger)
	//initialize middlewares
	mw := middleware.NewMiddleware(logger, queries)

	//initialize handlers
	hotelHandler := hh.NewHotelHandler(logger, hotelService)
	roomHandler := rh.NewRoomHandler(logger, roomService)
	userHandler := uh.NewUserHandler(logger, userService)

	//register routes

	userRoutes := []route{
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
			handler: userHandler.Login,
			middlewares: []gin.HandlerFunc{
				mw.Authorize(),
			},
		},
	}
	roomRoutes := []route{
		{
			path:    "/reserve",
			method:  http.MethodPost,
			handler: roomHandler.Reserve,
		},
	}
	hotelRoutes := []route{
		{
			path:    "/register",
			method:  http.MethodPost,
			handler: hotelHandler.Register,
			middlewares: []gin.HandlerFunc{
				mw.Authorize(),
			},
		},
	}

	allRoutes := append(userRoutes, append(hotelRoutes, roomRoutes...)...)

	r := gin.Default()

	RegisterRoutes(&r.RouterGroup, allRoutes)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
	log.Println("server started successfully")

}

func RegisterRoutes(g *gin.RouterGroup, routes []route) {
	for _, route := range routes {
		route.middlewares = append(route.middlewares, route.handler)
		g.Handle(route.method, route.path, route.middlewares...)
	}
}
