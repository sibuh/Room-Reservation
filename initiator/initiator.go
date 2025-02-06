package initiator

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"reservation/internal/service/hotel"
	"reservation/internal/service/payment"
	"reservation/internal/service/room"
	roomtype "reservation/internal/service/room_type"
	"reservation/internal/service/user"

	hh "reservation/internal/handler/hotel"
	"reservation/internal/handler/middleware"
	pmt "reservation/internal/handler/payment"
	rh "reservation/internal/handler/room"
	rth "reservation/internal/handler/room_type"
	uh "reservation/internal/handler/user"

	"reservation/internal/storage/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"

	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type route struct {
	path        string
	method      string
	handler     gin.HandlerFunc
	middlewares []gin.HandlerFunc
}

// var hotelHandler hh.HotelHandler
// var roomHandler rh.RoomHandler
// var userHandler uh.UserHandler
// var mw middleware.Middleware
// var paymentHandler pmt.PaymentHandler

func Initiate() {

	//initialize viper
	viper.SetConfigFile("config/config.yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}

	//create connection to database
	connString := viper.GetString("db.conn")
	pool, err := pgxpool.NewWithConfig(context.Background(), CreateDBConfig(connString))
	if err != nil {
		log.Fatal("failed to create connection pool ", err)
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Fatal("failed to acquire connection from pool", err)
	}

	//create logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	//create database
	_, err = conn.Exec(context.Background(), fmt.Sprintf("create database if not exists %s", viper.GetString("db.dbname")))
	if err != nil {
		log.Fatal("failed to create database", err)
	}
	//do database migration
	DoMigration(connString, "internal/storage/schema")

	// initialize storage layer
	queries := db.New(conn)
	//load env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	token_key := os.Getenv("TOKEN_KEY")
	stripePublishableKey := os.Getenv("STRIPE_PUBLISHABLE_KEY")
	paypal_client_id := os.Getenv("PAYPAL_CLIENT_ID")
	paypal_client_secret := os.Getenv("PAYPAL_CLIENT_SECRET")
	duration := viper.GetDuration("token.expire_after")
	cancellationTime := viper.GetDuration("reservation.cancellation_time")
	stripeSecretKey := os.Getenv("STRIPE_SECRET_KEY")

	//initialize services
	userService := user.NewUserService(logger, queries, token_key, duration)
	roomService := room.NewRoomService(pool, queries, logger, cancellationTime)
	hotelService := hotel.NewHotelService(queries, logger, pool)
	paymentService := payment.NewPaymentService(logger, queries, payment.PaymentProviderConfig{
		BaseURL:      viper.GetString("paypal.base_url"),
		ReturnURL:    viper.GetString("paypal_return_url"),
		CancelURL:    viper.GetString("paypal_return_url"),
		ClientID:     paypal_client_id,
		ClientSecret: paypal_client_secret,
		StripeSecret: stripeSecretKey,
	})
	roomTypeService := roomtype.NewRoomTypeService(logger, queries)

	//initialize casbin policy enforcer
	e := CasbinEnforcer("casbin/casbin.conf", "casbin/policy.csv")

	//initialize middlewares
	mw := middleware.NewMiddleware(logger, queries, token_key, e)

	//initialize handlers
	hotelHandler := hh.NewHotelHandler(logger, hotelService)
	roomHandler := rh.NewRoomHandler(logger, roomService)
	userHandler := uh.NewUserHandler(logger, userService)
	paymentHandler := pmt.NewPaymentHandler(logger, paymentService, stripePublishableKey)
	roomTypeHandler := rth.NewRoomTypeHandler(logger, roomTypeService)
	r := gin.Default()
	// r.Static("public", "./public")
	gin.SetMode(gin.ReleaseMode)
	gin.ForceConsoleColor()

	//register error handler for all routes
	r.Use(middleware.ErrorHandler())

	allRoutes := ListRoutes(roomHandler, hotelHandler, userHandler, paymentHandler, roomTypeHandler, mw)
	for _, rg := range allRoutes {
		RegisterRoutes(&r.RouterGroup, rg)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000"
	}
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
	log.Println("server started successfully")

}
