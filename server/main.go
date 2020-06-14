package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	blogpb "github.com/snowpato/sample-docker-restful-with-go/proto"
	"google.golang.org/grpc"
)

var goPort int

func init() {
	flag.IntVar(&goPort, "port", 50051, "Defines the local port to run the server.")
	flag.IntVar(&dbConn.Port, "mport", 27017, "Defines the mongodb listening port.")
	flag.StringVar(&dbConn.Host, "mongo", "localhost", "Defines the mongodb listening host.")
}

func main() {
	// Configure 'log' package to give file name and line number on eg. log.Fatal
	// Pipe flags to one another (log.LstdFLags = log.Ldate | log.Ltime)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Printf("Starting server on port :%d...\n", goPort)

	// Start our listener, 50051 is the default gRPC port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", goPort))
	// Handle errors if any
	if err != nil {
		log.Fatalf("Unable to listen on port :%d: %v", goPort, err)
	}

	//Connect to the Database
	err = dbConn.Connect()
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v\n", err)
	} else {
		fmt.Println("Connected to Mongodb")
	}

	//Set the collection through the database.
	dbConn.SetCollection("mydb", "blog")

	// Set options, here we can configure things like TLS support
	opts := []grpc.ServerOption{}
	// Create new gRPC server with (blank) options
	s := grpc.NewServer(opts...)
	// Create BlogService type
	srv := &BlogServiceServer{}
	// Register the service with the server
	blogpb.RegisterBlogServiceServer(s, srv)

	// Start the server in a child routine
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Printf("Server succesfully started on port :%d\n", goPort)

	// Right way to stop the server using a SHUTDOWN HOOK
	// Create a channel to receive OS signals
	c := make(chan os.Signal)

	// Relay os.Interrupt to our channel (os.Interrupt = CTRL+C)
	// Ignore other incoming signals
	signal.Notify(c, os.Interrupt)

	// Block main routine until a signal is received
	// As long as user doesn't press CTRL+C a message is not passed and our main routine keeps running
	<-c

	// After receiving CTRL+C Properly stop the server
	fmt.Println("\nStopping the server...")
	s.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB connection")
	dbConn.Disconnect()
	fmt.Println("Done.")
}
