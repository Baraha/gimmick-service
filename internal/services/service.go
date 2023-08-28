package services

import (
	"context"
	"google.golang.org/grpc"
	"mime/multipart"
)

// PortRepository is a port repository for the port service
type PortRepository interface {
	SetNewGRPCService(ctx context.Context, fileName string) error
	SetNewProtoFile(ctx context.Context, file Ifile, fileName string) error
}

type Ifile interface {
	Open() (multipart.File, error)
}

type IFileStorage interface {
	SetFile(fileFormat Ifile, fileName string) error
	DefaultDir() string
	DeleteFile(fileName string) error
}

type IGRPCService interface {
	LoadSpec(protoFileName string) error
	Handler(srv interface{}, ctx context.Context, dec func(interface{}) error,
		interceptor grpc.UnaryServerInterceptor) (interface{}, error)
}

// PortService is a port service
type PortService struct {
	dynGRPCHandler IGRPCService
	fs             IFileStorage
}

// NewPortService creates a new port service
func NewPortService(fs IFileStorage, dynGrpcHandler IGRPCService) PortRepository {
	return PortService{
		dynGRPCHandler: dynGrpcHandler,
		fs:             fs,
	}
}

// SetNewGRPCService set new grpc service in file storage
func (s PortService) SetNewGRPCService(ctx context.Context, fileName string) error {
	return s.dynGRPCHandler.LoadSpec(fileName)
}

func (s PortService) SetNewProtoFile(ctx context.Context, file Ifile, fileName string) error {
	return s.fs.SetFile(file, fileName)
}
