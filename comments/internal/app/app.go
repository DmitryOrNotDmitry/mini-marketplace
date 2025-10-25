package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"route256/comments/internal/handler"
	"route256/comments/internal/infra/config"
	"route256/comments/internal/infra/repository/postgres"
	"route256/comments/internal/service"
	"route256/comments/pkg/api/comments/v1"

	"route256/cart/pkg/logger"
	"route256/cart/pkg/myerrgroup"
	"route256/loms/pkg/grpc/interceptor"
	"route256/loms/pkg/http/middleware"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// App создает компоненты для сервиса comments
type App struct {
	Config       *config.Config
	grpcServer   *grpc.Server
	grpcGWServer *http.Server
}

// NewApp конструктор главного приложения.
func NewApp(ctx context.Context, configPath string) (*App, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{Config: c}

	app.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.Logging,
			interceptor.Validate,
		),
	)

	reflection.Register(app.grpcServer)

	shards := make([]*postgres.Shard, 0, len(app.Config.DB.Shards))
	for _, shardConfig := range app.Config.DB.Shards {
		postgresDSN := dsnBuilder(shardConfig.User, shardConfig.Password, shardConfig.Host,
			shardConfig.Port, shardConfig.DBName)

		shardPool, poolErr := newPool(ctx, postgresDSN)
		if poolErr != nil {
			return nil, fmt.Errorf("newPool: %w", poolErr)
		}

		shards = append(shards, &postgres.Shard{
			Pool:           shardPool,
			BucketPosition: shardConfig.BucketPosition,
		})
	}

	buckets := app.Config.DB.Buckets
	shardManager, err := postgres.NewShardManager(postgres.GetMurmur3Hashing(buckets), buckets, shards)
	if err != nil {
		return nil, fmt.Errorf("postgres.NewShardManager: %w", err)
	}

	commentRepository := postgres.NewCommentRepository(shardManager)

	duration, err := time.ParseDuration("1s")
	if err != nil {
		return nil, fmt.Errorf("time.ParseDuration: %w", err)
	}
	commentService := service.NewCommentService(commentRepository, duration)

	commentsHandler := handler.NewCommentServerGRPC(commentService)

	comments.RegisterCommentsServiceV1Server(app.grpcServer, commentsHandler)

	return app, nil
}

func dsnBuilder(user, password, host string, port int64, dbname string) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)
}

func newPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig (dsn=%s): %w", dsn, err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig (dsn=%s): %w", dsn, err)
	}

	return pool, nil
}

// ListenAndServeGRPC запускает gRPC-сервер приложения.
func (a *App) ListenAndServeGRPC() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", a.Config.Server.GRPCPort))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	logger.Infow(fmt.Sprintf("Comments service listening gRPC at port %s", a.Config.Server.GRPCPort))

	return a.grpcServer.Serve(listener)
}

// ListenAndServeGRPCGateway запускает gRPC-gateway для gRPC-сервера приложений.
func (a *App) ListenAndServeGRPCGateway(ctx context.Context) error {
	conn, err := grpc.NewClient(
		fmt.Sprintf(":%s", a.Config.Server.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("grpc.NewClient: %w", err)
	}

	gwMux := runtime.NewServeMux()

	err = comments.RegisterCommentsServiceV1Handler(ctx, gwMux, conn)
	if err != nil {
		return fmt.Errorf("comments.RegisterCommentsServiceV1Handler: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwMux)

	handler := middleware.CORSAllPass(mux)

	readHeaderTimeout, err := time.ParseDuration(a.Config.Server.GRPCGateWay.ReadHeaderTimeout)
	if err != nil {
		return fmt.Errorf("time.ParseDuration: %w", err)
	}
	writeTimeout, err := time.ParseDuration(a.Config.Server.GRPCGateWay.WriteTimeout)
	if err != nil {
		return fmt.Errorf("time.ParseDuration: %w", err)
	}
	idleTimeout, err := time.ParseDuration(a.Config.Server.GRPCGateWay.IdleTimeout)
	if err != nil {
		return fmt.Errorf("time.ParseDuration: %w", err)
	}

	a.grpcGWServer = &http.Server{
		Addr:              fmt.Sprintf(":%s", a.Config.Server.HTTPPort),
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	logger.Infow(fmt.Sprintf("Comments service listening gRPC-Gateway (REST) at port %s", a.Config.Server.HTTPPort))

	return a.grpcGWServer.ListenAndServe()
}

// Shutdown gracefully останавливает приложение.
func (a *App) Shutdown(ctx context.Context) error {
	errGroup, ctx := myerrgroup.WithContext(ctx)

	errGroup.Go(func() error {
		return a.grpcGWServer.Shutdown(ctx)
	})

	errGroup.Go(func() error {
		successGraceful := make(chan struct{})
		go func() {
			a.grpcServer.GracefulStop()
			close(successGraceful)
		}()

		select {
		case <-ctx.Done():
			a.grpcServer.Stop()
			return ctx.Err()
		case <-successGraceful:
			return nil
		}
	})

	return errGroup.Wait()
}
