package presence

import (
	"GameApp/contract/golang/presence"
	"GameApp/param"
	"GameApp/pkg/protobuf"
	"GameApp/pkg/slice"
	"context"
	"google.golang.org/grpc"
)

type Client struct {
	client presence.PresenceServiceClient
}

func New(conn *grpc.ClientConn) Client {
	return Client{presence.NewPresenceServiceClient(conn)}
}

func (c Client) GetPresence(ctx context.Context, request param.GetPresenceRequest) (param.GetPresenceResponse, error) {
	res, err := c.client.GetPresence(ctx, &presence.GetPresenceRequest{UserIds: slice.MapFromUintToUint64(request.UserIDs)})
	if err != nil {
		return param.GetPresenceResponse{}, err
	}
	return protobuf.MapGetPresenceResponseFromProtobuf(res), nil
}
