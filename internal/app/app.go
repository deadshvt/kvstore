package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"

	"github.com/deadshvt/kvstore/config"
	pairHandler "github.com/deadshvt/kvstore/internal/delivery/http/pair"
	userHandler "github.com/deadshvt/kvstore/internal/delivery/http/user"
	"github.com/deadshvt/kvstore/internal/middleware"
	"github.com/deadshvt/kvstore/internal/middleware/monitoring"
	"github.com/deadshvt/kvstore/internal/middleware/protection"
	pairRepo "github.com/deadshvt/kvstore/internal/repository/pair"
	pairDB "github.com/deadshvt/kvstore/internal/repository/pair/database"
	userRepo "github.com/deadshvt/kvstore/internal/repository/user"
	userDB "github.com/deadshvt/kvstore/internal/repository/user/database"
	"github.com/deadshvt/kvstore/internal/security"
	pairUsecase "github.com/deadshvt/kvstore/internal/usecase/pair"
	userUsecase "github.com/deadshvt/kvstore/internal/usecase/user"
	"github.com/deadshvt/kvstore/pkg/logger"
)

func Run() {
	baseLogger, err := logger.Init("kvstore")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	lg := logger.NewLogger(baseLogger, "kvstore")

	lg.Info().Msg("Initialized logger")

	config.Load(".env")

	lg.Info().Msg("Loaded .env file")

	ctx := context.Background()

	// User

	// JWT
	jwtService := security.NewJWTService(jwt.SigningMethodHS256, os.Getenv("JWT_SECRET"))

	// DB
	uDB, err := userDB.NewUserDB(ctx, os.Getenv("USER_DB_TYPE"))
	if err != nil {
		lg.Fatal().Msgf("Failed to connect to user database: %v", err)
	}
	defer func() {
		err = uDB.Disconnect(ctx)
		if err != nil {
			lg.Error().Msgf("Failed to disconnect from user database: %v", err)
		}
	}()

	lg.Info().Msg("Connected to user database")

	// Repository
	uRepoLogger := logger.NewLogger(baseLogger, "user-repository")
	uRepo := userRepo.NewRepository(uDB, uRepoLogger)

	// Usecase
	uUsecaseLogger := logger.NewLogger(baseLogger, "user-usecase")
	uUsecase := userUsecase.NewUsecase(uRepo, jwtService, os.Getenv("USER_ENCRYPTION_KEY"), uUsecaseLogger)

	// Handler
	uHandlerLogger := logger.NewLogger(baseLogger, "user-handler")
	uHandler := userHandler.NewHandler(uUsecase, uHandlerLogger)

	// Pair

	// DB
	pDB, err := pairDB.NewPairDB(ctx, os.Getenv("PAIR_DB_TYPE"))
	if err != nil {
		lg.Fatal().Msgf("Failed to connect to pair database: %v", err)
	}
	defer func() {
		err = pDB.Disconnect(ctx)
		if err != nil {
			lg.Error().Msgf("Failed to disconnect from pair database: %v", err)
		}
	}()

	lg.Info().Msg("Connected to pair database")

	// Repository
	pRepoLogger := logger.NewLogger(baseLogger, "pair-repository")
	pRepo := pairRepo.NewRepository(pDB, pRepoLogger)

	// Usecase
	pUsecaseLogger := logger.NewLogger(baseLogger, "pair-usecase")
	pUsecase := pairUsecase.NewUsecase(pRepo, os.Getenv("PAIR_ENCRYPTION_KEY"), pUsecaseLogger)

	// Handler
	pHandlerLogger := logger.NewLogger(baseLogger, "pair-handler")
	pHandler := pairHandler.NewHandler(pUsecase, pHandlerLogger)

	// Middleware

	// Monitoring
	lgr := monitoring.NewLgr(logger.NewLogger(baseLogger, "lgr"))

	metricsCollector := monitoring.NewMetricsCollector(logger.NewLogger(baseLogger, "metricsCollector"))
	metricsCollector.Register()

	requestIDGenerator := monitoring.NewRequestIDGenerator(logger.NewLogger(baseLogger, "requestIDGenerator"))

	// Security
	authenticator := protection.NewAuthenticator(jwtService, logger.NewLogger(baseLogger, "authenticator"))

	content := protection.NewContent("application/json", logger.NewLogger(baseLogger, "content"))

	limiter := rate.NewLimiter(1, 10)
	rateLimiter := protection.NewRateLimiter(limiter, logger.NewLogger(baseLogger, "rateLimiter"))

	recoverer := protection.NewRecoverer(logger.NewLogger(baseLogger, "recoverer"))

	publicChain := middleware.ChainMiddleware(
		recoverer.Middleware,
		requestIDGenerator.Middleware,
		lgr.Middleware,
		metricsCollector.Middleware,
		content.Middleware,
		rateLimiter.Middleware,
	)

	privateChain := middleware.ChainMiddleware(
		recoverer.Middleware,
		requestIDGenerator.Middleware,
		lgr.Middleware,
		metricsCollector.Middleware,
		content.Middleware,
		rateLimiter.Middleware,
		authenticator.Middleware,
	)

	// Router
	r := mux.NewRouter()

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.Use(publicChain)

	apiRouter.HandleFunc("/login", uHandler.Login).Methods(http.MethodPost)

	apiAuthRouter := r.PathPrefix("/api").Subrouter()
	apiAuthRouter.Use(privateChain)

	apiAuthRouter.HandleFunc("/write", pHandler.SetPairs).Methods(http.MethodPost)
	apiAuthRouter.HandleFunc("/read", pHandler.GetPairs).Methods(http.MethodPost)

	// Server
	srv := &http.Server{
		Addr:    ":" + os.Getenv("SERVER_PORT"),
		Handler: r,
	}

	lg.Info().Msgf("Starting server on port :%s", os.Getenv("SERVER_PORT"))

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Error().Msgf("Failed to start server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	lg.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		lg.Error().Msgf("Failed to shutdown http server: %v", err)
	}

	lg.Info().Msg("Server gracefully stopped")
}
