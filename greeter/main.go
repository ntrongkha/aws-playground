package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lotusirous/greeter/config"
)

// DB defines a serialization for storing secret manager.
// The data from secretsmanager -> json format.
// sample data: {"username":"kha","password":"1","engine":"postgres","host":"127.0.0.1","port":"5432","dbname":"postgres"}
type DB struct {
	Username string `json:"username"`
	Engine   string `json:"engine"`
	Password string `json:"password"`
}

type worker struct {
	invokeCount int
	DB          DB
}

// HandleEvent performs an event handling.
// It follows the lambda handler signature.
func (s *worker) HandleEvent(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("HandleEvent time: %d\n", s.invokeCount)
	s.invokeCount++
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("time %d read from secrets: %v", s.invokeCount, s.DB),
	}, nil
}

func main() {
	log.Println("cold start start")
	var (
		ctx      = context.Background()
		secretID = os.Getenv("SECRET_KEY")
	)
	var db DB
	err := config.UnmarshalSecret(ctx, secretID, &db)
	if err != nil {
		log.Fatalf("unable to read secret: %v\n", err)
	}
	w := worker{DB: db}
	// cold start stop
	log.Println("cold start stop")

	lambda.Start(w.HandleEvent) // "Warm Start" is calling from the handler
}
