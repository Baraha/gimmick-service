package proto_generic

import (
	"context"
	"github.com/Baraha/gimmick-service/internal/services"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
)

type GRPCService struct {
	sdMap      map[string]*desc.ServiceDescriptor
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewGRPCService(grpcServer *grpc.Server, listener net.Listener) services.IGRPCService {
	return &GRPCService{grpcServer: grpcServer, listener: listener,
		sdMap: make(map[string]*desc.ServiceDescriptor, 24)}
}

// Parse protofile, create grpc.ServiceDesc, register
func (s *GRPCService) LoadSpec(protoFileName string) error {
	p := protoparse.Parser{}
	fdlist, err := p.ParseFiles(protoFileName)
	if err != nil {
		return err
	}
	for _, fd := range fdlist {
		for _, rsd := range fd.GetServices() {

			if _, exist := s.sdMap[rsd.GetFullyQualifiedName()]; exist {
				return errors.New("service already exist")
			}
			s.sdMap[rsd.GetFullyQualifiedName()] = rsd
			gsd := grpc.ServiceDesc{ServiceName: rsd.GetFullyQualifiedName(), HandlerType: (*interface{})(nil)}
			for _, m := range rsd.GetMethods() {
				gsd.Methods = append(gsd.Methods, grpc.MethodDesc{MethodName: m.GetName(), Handler: s.Handler})
			}
			s.grpcServer.GracefulStop()
			s.grpcServer = grpc.NewServer()
			s.grpcServer.RegisterService(&gsd, s.grpcServer)

			go func() {
				if err := s.grpcServer.Serve(s.listener); err != nil {
					log.Printf("Error : %v", err)
				}
			}()

		}
	}
	return nil
}

func (s *GRPCService) Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	stream := grpc.ServerTransportStreamFromContext(ctx)
	arr := strings.Split(stream.Method(), "/")
	serviceName := arr[1]
	methodName := arr[2]
	service := s.sdMap[serviceName]
	method := service.FindMethodByName(methodName)
	input := dynamic.NewMessage(method.GetInputType())

	err := dec(input)
	if err != nil {
		return nil, err
	}
	jsonInput, err := input.MarshalJSON()
	log.Printf("Input:%s Err:%v \n", jsonInput, err)
	//jsonOutput:=invokeServiceViaReflectionOrHttp(jsonInput)
	jsonOutput := `{"message":"response"}`

	output := dynamic.NewMessage(method.GetOutputType())
	output.UnmarshalJSON([]byte(jsonOutput))
	return output, nil
}
