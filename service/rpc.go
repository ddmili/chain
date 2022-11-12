package service

import (
	"chain/config"
	"chain/internal/log"
	"net"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GrpcServer rpc service
func GrpcServer(bc *ChainService, lg log.Logger, cg *config.Config) error {
	lis, err := net.Listen("tcp", cg.RPCConfig.Port)
	if err != nil {
		log.Panic("start grpc service failed:%s", err)
		return err
	}
	// 创建 RPC 服务容器
	grpcServer := grpc.NewServer()

	RegisterChainServer(grpcServer, bc)

	reflection.Register(grpcServer)

	log.Info("start grpc service port:%s", cg.RPCConfig.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Panic("start grpc service failed:%s", err)
		return err
	}
	return nil
}
