package interceptor

import (
	"context"
	"gophkeeper/internal/identity"
	"slices"

	pb "gophkeeper/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ignoreMethods = []string{pb.Keeper_Login_FullMethodName, pb.Keeper_Registration_FullMethodName}
	Token         = "token"
	UserId        = "user_id"
)

type authInterceptor struct {
	identityProvider identity.IdentityProvider
}

func GetAuthInterceptor(provider identity.IdentityProvider) authInterceptor {
	return authInterceptor{identityProvider: provider}
}

func (a *authInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if !slices.Contains(ignoreMethods, info.FullMethod) {
		var token string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("token")
			if len(values) > 0 {
				token = values[0]
			}
		}
		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}
		userId, err := a.identityProvider.ParseToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
		}
		md := metadata.New(map[string]string{"user_id": string(userId)})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return handler(ctx, req)
}
