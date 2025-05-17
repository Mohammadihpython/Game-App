package presence

import (
	"GameApp/contract/goproto/presence"
	"GameApp/param"
	"GameApp/pkg/protobufMapper"
	"GameApp/pkg/slice"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	address string
}

func New(address string) Client {
	return Client{address: address}
}

func (c Client) GetPresence(ctx context.Context, request param.GetPresenceRequest) (param.GetPresenceResponse, error) {

	conn, err := grpc.NewClient(c.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return param.GetPresenceResponse{}, err
	}
	defer conn.Close()
	client := presence.NewPresenceServiceClient(conn)
	res, err := client.GetPresence(ctx, &presence.GetPresenceRequest{UserIds: slice.MapFromUintToUint64(request.UserIDs)})
	if err != nil {
		return param.GetPresenceResponse{}, err
	}
	return protobufMapper.MapGetPresenceResponseFromProtobuf(res), nil
}

func (c Client) Upsert(ctx context.Context, request param.UpsertPresenceRequest) (param.UpsertPresenceResponse, error) {
	conn, err := grpc.NewClient(c.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return param.UpsertPresenceResponse{}, err
	}
	defer conn.Close()
	client := presence.NewPresenceServiceClient(conn)
	_, err = client.UpsertPresence(ctx, &presence.UpsertPresenceRequest{
		UserId:    uint64(request.UserID),
		Timestamp: request.Timestamp,
	})
	if err != nil {
		return param.UpsertPresenceResponse{}, err
	}
	return param.UpsertPresenceResponse{}, nil
}
