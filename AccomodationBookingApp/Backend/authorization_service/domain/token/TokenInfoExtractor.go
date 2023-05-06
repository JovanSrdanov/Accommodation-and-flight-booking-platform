package token

import (
	"context"
	"fmt"
	"github.com/o1egl/paseto"
	"google.golang.org/grpc/metadata"
	"strings"
)

func ExtractInfoFromToken(ctx context.Context, infoType string) (interface{}, error) {
	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata provided")
	}

	values := metaData["authorization"]
	token := strings.TrimPrefix(values[0], "Bearer ")
	var footerData map[string]interface{}
	if err := paseto.ParseFooter(token, &footerData); err != nil {
		return nil, fmt.Errorf("failed to parse token footer")
	}

	return footerData[infoType], nil
}
