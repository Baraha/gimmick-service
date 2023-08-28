package main

import (
	"github.com/Baraha/gimmick-service/internal/config"
	"github.com/Baraha/gimmick-service/internal/repository/filestorage"
	proto_generic "github.com/Baraha/gimmick-service/internal/repository/proto-generic"
	"github.com/Baraha/gimmick-service/internal/services"
	dyn_grpc_transport "github.com/Baraha/gimmick-service/internal/transport/http/dyn-grpc-transport"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	// read config from env
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	fs := filestorage.NewFileStorage(*cfg)
	listenErr := make(chan error, 1) // add error handler

	grpcServer, grpcListener := initKitGRPCServer(cfg, listenErr)
	defer func() {
		err = grpcListener.Close()
		if err != nil {
			// do we really need it? i think no because we already closed it by cmux.Close()
			// netLogger.Warn().Err(err).Msgf("failed to close net.Listen %+w - %+v", err, err)
		}
	}()

	dynGRPCs := proto_generic.NewGRPCService(grpcServer, grpcListener)

	log.Printf("init services")
	portService := services.NewPortService(fs, dynGRPCs)
	dGRPCTransport := dyn_grpc_transport.NewHttpDynGRPCServer(portService)

	r := router.New()
	r.POST("/SetNewGRPCService", dGRPCTransport.SetNewGRPCService)
	r.POST("/SetNewProtoFile", dGRPCTransport.SetNewProtoFile)

	// start HTTP server

	log.Printf("starting service on %s adress", cfg.HTTPServer.Address)
	if err := fasthttp.ListenAndServe(cfg.HTTPServer.Address, r.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe: %v", err)
	}

	log.Printf("Have a nice day!")

	return nil
}

func initKitGRPCServer(appConfig *config.Configuration, listenErr chan error) (*grpc.Server, net.Listener) {

	var grpcServer *grpc.Server

	grpcServer = grpc.NewServer()

	l, err := net.Listen(appConfig.GRPCServer.Network, appConfig.GRPCServer.Address)
	if err != nil {
		log.Fatalf("Error in initKitGRPCServer: %v", err)
	}
	time.Sleep(10 * time.Millisecond)
	return grpcServer, l
}

func RunGRPCServer(grpcSrv *grpc.Server, l net.Listener, listenErr chan error) {
	reflection.Register(grpcSrv)
	log.Printf("starting grpc server on %s", l.Addr())
	if err := grpcSrv.Serve(l); err != nil {
		listenErr <- err
	}
}
