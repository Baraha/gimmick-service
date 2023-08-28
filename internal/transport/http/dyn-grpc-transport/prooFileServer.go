package dyn_grpc_transport

import (
	"encoding/json"
	"github.com/Baraha/gimmick-service/internal/services"
	"github.com/Baraha/gimmick-service/internal/transport/http/errors"
	"github.com/Baraha/gimmick-service/internal/transport/http/models"
	"github.com/valyala/fasthttp"
	"log"
)

type httpDynGRPCServer struct {
	portDynGRPC services.PortRepository
}

func NewHttpDynGRPCServer(portDynGRPC services.PortRepository) httpDynGRPCServer {
	return httpDynGRPCServer{portDynGRPC}
}

func (service *httpDynGRPCServer) SetNewProtoFile(ctx *fasthttp.RequestCtx) {
	header, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(errors.ErrorFileNotFound.Error(), fasthttp.StatusBadRequest)
	}
	log.Println("save newGRPC File")
	if err := service.portDynGRPC.SetNewProtoFile(ctx, header, header.Filename); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
	}
}

func (service *httpDynGRPCServer) SetNewGRPCService(ctx *fasthttp.RequestCtx) {
	var dto models.GrpcServerInputDTO
	if err := json.Unmarshal(ctx.Request.Body(), &dto); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
	}
	if err := service.portDynGRPC.SetNewGRPCService(ctx, dto.ServerName); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
	}
}
