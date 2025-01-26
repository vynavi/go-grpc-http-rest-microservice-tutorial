package server
package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"

	
	"github.com/vynavi/go-grpc-http-rest-microservice-tutorial/pkg/protocol/grpc"
	"github.com/vynavi/go-grpc-http-rest-microservice-tutorial/pkg/service/v1"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	GRPCPort string

	// DB Datastore parameters section
	DatastoreDBHost     string
	DatastoreDBUser     string
	DatastoreDBPassword string
	DatastoreDBSchema   string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// Get configuration from command-line flags
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	// Validate gRPC port
	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	// Create MySQL DSN string for connecting to the database
	param := "parseTime=true"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBSchema,
		param)

	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	// Validate the database connection with Ping()
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Log successful DB connection
	log.Println("Successfully connected to the database")

	// Create ToDoServiceServer instance
	v1API := v1.NewToDoServiceServer(db)

	// Run the gRPC server
	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}

func main() {
	// Run the server
	if err := RunServer(); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
