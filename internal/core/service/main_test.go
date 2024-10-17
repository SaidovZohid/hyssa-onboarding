package service_test

import (
	"context"
	"learning/hyssa-learn/internal/config"
	"learning/hyssa-learn/internal/core/repository"
	"learning/hyssa-learn/internal/pkg/logger"
	"learning/hyssa-learn/internal/transport/grpc"
	"learning/hyssa-learn/internal/transport/handlers"
	"learning/hyssa-learn/migrate"
	"learning/hyssa-learn/pkg/container"
	"log"
	"net"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	psqlContainer *container.PostgreSQLContainer
	server        *httptest.Server
	cfg           *config.Config
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	s.cfg = config.Load()

	logger.SetLogger(&s.cfg.Logging)

	psqlContainer, err := container.NewPostgreSQLContainer(ctx)
	s.NoError(err)

	s.psqlContainer = psqlContainer

	err = migrate.Migrate(psqlContainer.GetDSN(), migrate.Migrations)
	s.NoError(err)

	repos := repository.New(ctx, psqlContainer.GetDSN())

	grpcServer := grpc.New(repos, s.cfg)

	lis, err := net.Listen("tcp", s.cfg.Grpc.URL)
	if err != nil {
		log.Fatal("Error while listening: ", err)
		return
	}

	go func() {
		logger.Log.Info("starting grpc server on " + s.cfg.Grpc.URL)
		grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	srv := httptest.NewServer(handlers.New(ctx, s.cfg))

	s.server = srv
}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.psqlContainer.Terminate(ctx))

	s.server.Close()
}
