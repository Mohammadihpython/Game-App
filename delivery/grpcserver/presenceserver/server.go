package presenceserver

import (
	"GameApp/contract/goproto/presence"
	"GameApp/param"
	"GameApp/pkg/protobufMapper"
	"GameApp/pkg/slice"
	"GameApp/service/presenceservice"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

type Server struct {
	presence.UnimplementedPresenceServiceServer
	svc presenceservice.Service
}

func New(svc presenceservice.Service) Server {
	return Server{UnimplementedPresenceServiceServer: presence.UnimplementedPresenceServiceServer{}, svc: svc}
}

func (s Server) GetPresence(ctx context.Context, req *presence.GetPresenceRequest) (*presence.GetPresenceResponse, error) {

	res, err := s.svc.GetPresence(ctx, param.GetPresenceRequest{UserIDs: slice.MapFromUint64ToUint(req.GetUserIds())})
	if err != nil {
		return nil, err
	}
	fmt.Println(res)

	return protobufMapper.MapGetPresenceResponseToProtobuf(res), nil

}

func (s Server) UpsertPresence(ctx context.Context, req *presence.UpsertPresenceRequest) (*emptypb.Empty, error) {

	res, err := s.svc.Upsert(ctx, param.UpsertPresenceRequest{UserID: uint(req.GetUserId()), Timestamp: req.GetTimestamp()})
	if err != nil {
		return nil, err
	}
	fmt.Println(res)

	return &emptypb.Empty{}, nil

}

func (s Server) Start() {
	// listener := tcp port
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", 8070))
	if err != nil {
		panic(err)
	}

	//grpc server
	grpcServer := grpc.NewServer()

	// pb presence server register into grpc server
	presence.RegisterPresenceServiceServer(grpcServer, &s)
	//	server grpcServer by listener
	if err := grpcServer.Serve(listner); err != nil {
		log.Fatal("Could not start gRPC server")
	}
}
