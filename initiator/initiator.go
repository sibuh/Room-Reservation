package initiator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type route struct {
	path        string
	method      string
	handler     gin.HandlerFunc
	middlewares []gin.HandlerFunc
}

func Initiate() {

	//initialize viper
	viper.SetConfigFile("config/config.yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}

	//create connection to database
	connString := viper.GetString("db_conn")
	pool, err := pgxpool.NewWithConfig(context.Background(), CreateDBConfig(connString))
	if err != nil {
		log.Fatal("failed to create connection pool", err)
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Fatal("failed to create connection from pool", err)
	}

	//create logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	_, err = conn.Exec(context.Background(), "create database if not exists reservation")
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
	key := os.Getenv("TOKEN_KEY")
	duration := viper.GetDuration("token_duration")

	//initialize services
	userService := user.NewUserService(logger, queries, key, duration)
	roomService := room.NewRoomService(queries, "url")
	hotelService := hotel.NewHotelService(queries, logger)

	//initialize middlewares
	mw := middleware.NewMiddleware(logger, queries, key)

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
			handler: userHandler.Refresh,
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
	gin.SetMode(gin.ReleaseMode)

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

func CreateDBConfig(url string) *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("connection released!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database!!")
	}

	return dbConfig
}

func DoMigration(connString, filePath string) {
	m, err := migrate.New(fmt.Sprintf("file://%s", filePath),
		"cockroachdb://"+strings.Split(connString, "//")[1])
	if err != nil {
		log.Fatal("failed to create migration instance", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal("failed to do migration: ", err)
		}
	}
}
