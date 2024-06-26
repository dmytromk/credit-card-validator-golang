package main

import (
	"card_validator/pb"
	"card_validator/validation"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedValidationServer
}

func (s *server) ValidateCard(ctx context.Context, in *pb.CardRequest) (*pb.ValidationResponse, error) {
	response := &pb.ValidationResponse{Valid: true}

	trimmedCardNumber := validation.RemoveWhitespace(&in.CardNumber)

	if digit := validation.IsDigit(&trimmedCardNumber); digit != nil {
		response.Valid = false
		response.Errors = append(response.Errors, digit)
	} else if luhn := validation.LuhnCheck(&trimmedCardNumber); luhn != nil {
		response.Valid = false
		response.Errors = append(response.Errors, luhn)
	} else if issuer := validation.IssuerCheck(&trimmedCardNumber); issuer != nil {
		response.Valid = false
		response.Errors = append(response.Errors, issuer)
	}

	if expiration := validation.Expiration(&in.ExpirationMonth, &in.ExpirationYear); expiration != nil {
		response.Valid = false
		response.Errors = append(response.Errors, expiration)
	}

	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen on port 5001: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterValidationServer(s, &server{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
