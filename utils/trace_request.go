package utils

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func GenerateRequestID() string {
	uuid1 := uuid.New()
	return fmt.Sprintf("%s", uuid1)
}

func GetRequestID(ctx context.Context) string {
	return ctx.Value("request_id").(string)
}

func WrapRequestID(reqID string) context.Context {
	ctx := context.WithValue(context.Background(), "request_id", reqID)
	return ctx
}
