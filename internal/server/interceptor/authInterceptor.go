package interceptor

import (
	"context"
	"gophkeeper/internal/identity"
	"slices"
	"strconv"

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

func (a *authInterceptor) UnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if !slices.Contains(ignoreMethods, info.FullMethod) {
		var token string
		var md metadata.MD
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
			return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
		}

		md = metadata.New(map[string]string{"user_id": strconv.Itoa(userId)})
		ctx = metadata.NewIncomingContext(ctx, md)
		return handler(ctx, req)
	}
	return handler(ctx, req)
}

func (a *authInterceptor) StreamingAuthInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if !slices.Contains(ignoreMethods, info.FullMethod) {
		var token string
		var md metadata.MD
		if md, ok := metadata.FromIncomingContext(ss.Context()); ok {
			values := md.Get("token")
			if len(values) > 0 {
				token = values[0]
			}
		}
		if len(token) == 0 {
			return status.Error(codes.Unauthenticated, "missing token")
		}
		userId, err := a.identityProvider.ParseToken(token)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "unauthenticated")
		}

		md = metadata.New(map[string]string{"user_id": strconv.Itoa(userId)})
		ctx := metadata.NewIncomingContext(ss.Context(), md)
		return handler(srv, newWrappedStream(ctx, ss))
	}

	return handler(srv, ss)
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

func newWrappedStream(ctx context.Context, s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s, ctx}
}
