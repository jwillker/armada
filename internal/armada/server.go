package armada

import (
	"github.com/G-Research/k8s-batch/internal/armada/api"
	"github.com/G-Research/k8s-batch/internal/armada/configuration"
	"github.com/G-Research/k8s-batch/internal/armada/repository"
	"github.com/G-Research/k8s-batch/internal/armada/server"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

func Serve(config *configuration.ArmadaConfig) (*grpc.Server, *sync.WaitGroup) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	grpcServer := grpc.NewServer()
	go func () {
		log.Printf("Grpc listening on %s", config.GrpcPort)
		defer log.Println("Stopping server.")

		db := redis.NewClient(&redis.Options{
			Addr:     config.Redis.Addr,
			Password: config.Redis.Password,
			DB:       config.Redis.Db,
		})

		jobRepository := &repository.RedisJobRepository{ Db: db }
		usageRepository := &repository.RedisUsageRepository{ Db: db }

		submitServer := &server.SubmitServer{ JobRepository: jobRepository }
		usageServer := &server.UsageServer { UsageRepository: usageRepository, PriorityHalfTime: time.Minute }
		aggregatedQueueServer := &server.AggregatedQueueServer{ JobRepository: jobRepository, UsageRepository: usageRepository }

		lis, err := net.Listen("tcp", config.GrpcPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		api.RegisterSubmitServer(grpcServer, submitServer)
		api.RegisterUsageServer(grpcServer, usageServer)
		api.RegisterAggregatedQueueServer(grpcServer, aggregatedQueueServer)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

		wg.Done()
	} ()
	return grpcServer, wg
}
