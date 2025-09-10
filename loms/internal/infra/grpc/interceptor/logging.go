package interceptor

import (
	"context"
	"fmt"
	"route256/cart/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)

	st := status.Convert(err)
	logger.Info(fmt.Sprintf("method=%s code=%d(%s) err_message='%s'", info.FullMethod, st.Code(), st.Code(), st.Message()))

	return resp, err
}
