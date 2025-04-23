package presenceserver

import (
	"GameApp/contract/golang/presence"
	"GameApp/param"
	"GameApp/pkg/protobuf"
	"GameApp/pkg/slice"
	"GameApp/service/presenceservice"
	"context"
	"fmt"
	"google.golang.org/grpc"
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

	return protobuf.MapGetPresenceResponseToProtobuf(res), nil

}

func (s Server) Start() {
	// listener := tcp port
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", 8070))
	if err != nil {
		panic(err)
	}
	// pbPresence server
	presenceSVCServer := Server{}
	//grpc server
	grpcServer := grpc.NewServer()

	// pbpresenceserver register into grpc server
	presence.RegisterPresenceServiceServer(grpcServer, &presenceSVCServer)
	//	server grpcServer by listener
	if err := grpcServer.Serve(listner); err != nil {
		log.Fatal("Could not start gRPC server")
	}
}
