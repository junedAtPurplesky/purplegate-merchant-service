package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/junedAtPurplesky/purplegate-merchant-service/internal/db"
	pb "github.com/junedAtPurplesky/purplegate-merchant-service/proto"
)

// MerchantServer implements the MerchantService gRPC server
type MerchantServer struct {
	pb.UnimplementedMerchantServiceServer
	repo *db.MerchantRepository
}

// NewMerchantServer creates a new instance of MerchantServer
func NewMerchantServer(repo *db.MerchantRepository) *MerchantServer {
	return &MerchantServer{
		repo: repo,
	}
}

// RegisterMerchant handles merchant registration
func (s *MerchantServer) RegisterMerchant(ctx context.Context, req *pb.RegisterMerchantRequest) (*pb.RegisterMerchantResponse, error) {
	// Validate input
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}
	if req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}

	// Create merchant in database
	merchant, err := s.repo.CreateMerchant(ctx, req.Name, req.Email)
	if err != nil {
		// Handle specific database errors
		if err.Error() == fmt.Sprintf("merchant with email %s already exists", req.Email) {
			return nil, status.Errorf(codes.AlreadyExists, "merchant with email %s already exists", req.Email)
		}
		return nil, status.Errorf(codes.Internal, "failed to create merchant: %v", err)
	}

	return &pb.RegisterMerchantResponse{
		MerchantId: merchant.MerchantId,
	}, nil
}

// GetMerchant retrieves merchant information by ID
func (s *MerchantServer) GetMerchant(ctx context.Context, req *pb.GetMerchantRequest) (*pb.GetMerchantResponse, error) {
	// Validate input
	if req.MerchantId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "merchant_id is required")
	}

	// Get merchant from database
	merchant, err := s.repo.GetMerchant(ctx, req.MerchantId)
	if err != nil {
		if err.Error() == fmt.Sprintf("merchant not found: %s", req.MerchantId) {
			return nil, status.Errorf(codes.NotFound, "merchant not found: %s", req.MerchantId)
		}
		return nil, status.Errorf(codes.Internal, "failed to get merchant: %v", err)
	}

	return merchant, nil
}

// StartServer starts the gRPC server
func (s *MerchantServer) StartServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMerchantServiceServer(grpcServer, s)

	fmt.Printf("Merchant service gRPC server listening on port %d\n", port)
	return grpcServer.Serve(lis)
}
