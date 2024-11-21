package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/algol-84/auth/internal/metric"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Инкрементим метрику
	metric.IncRequestCounter()

	log.Println("Metrics Interceptor")

	//	timeStart := time.Now()

	res, err := handler(ctx, req)
	// diffTime := time.Since(timeStart)

	// if err != nil {
	// 	metric.IncResponseCounter("error", info.FullMethod)
	// 	metric.HistogramResponseTimeObserve("error", diffTime.Seconds())
	// } else {
	// 	metric.IncResponseCounter("success", info.FullMethod)
	// 	metric.HistogramResponseTimeObserve("success", diffTime.Seconds())
	// }

	return res, err
}
