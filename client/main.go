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

// Package main implements a client for Greeter service.
package main

import (
	"context"

	"github.com/danztran/grpc_demo/api"
	"github.com/danztran/grpc_demo/util"
	"github.com/danztran/logger/log"
	"google.golang.org/grpc"
)

var (
	address = util.Getenv("SERVER_ADDR", "localhost:50051")
)

func main() {
	defer log.Infod()("execution time")

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewGreeterClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-util.StopSignal()
		log.Infof("terminating...")
		cancel()
	}()

	pool := util.NewWorkerPool(80)
	defer pool.Wait()

	// Contact the server and print out its response.
	for ctx.Err() == nil {
		pool.Run(func() {
			ctx := context.Background()
			r, err := c.SayHello(ctx, &api.HelloRequest{Name: "foo"})
			if err != nil {
				log.Warnf("could not greet: %v", err)
			}
			log.Infof("Greeting: %s", r.GetMessage())
		})
	}
}
