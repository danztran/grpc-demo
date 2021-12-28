/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/carousell/gologger/log"
	pb "github.com/danztran/grpc_demo/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
)

const (
	grpcAddr = ":50051"
	httpAddr = ":50052"
)

// Service is used to implement helloworld.GreeterServer.
type Service struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *Service) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if err := in.ValidateAll(); err != nil {
		return nil, err
	}
	time.Sleep(1000 * time.Millisecond)
	log.Infof("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// SayHello implements helloworld.GreeterServer
func (s *Service) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Infof("Received again: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName() + " again"}, nil
}

func main() {
	defer log.Infof("terminated!")
	ctx := context.Background()

	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	mux := runtime.NewServeMux()
	httpServer := http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}
	defer httpServer.Shutdown(ctx)

	svc := new(Service)
	pb.RegisterGreeterServer(grpcServer, svc)
	pb.RegisterGreeterHandlerServer(ctx, mux, svc)

	// close
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
		grpcServer.GracefulStop()
	}()

	errchan := make(chan error)

	// grpc
	go func() {
		log.Infof("grpc server starting at: %v", grpcAddr)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			errchan <- err
			return
		}
		if err := grpcServer.Serve(lis); err != nil {
			errchan <- err
			return
		}
	}()

	// http
	go func() {
		log.Infof("http server starting at: %v", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			errchan <- err
			return
		}
	}()

	select {
	case <-StopSignal():
	case err := <-errchan:
		log.Errorf("start error / %v", err)
	}

	log.Info("terminating...")
}

// StopSignal return a canceling signals channel (like INT, TERM)
func StopSignal() chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}
