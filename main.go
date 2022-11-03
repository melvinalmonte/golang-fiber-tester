package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
)

var fiberLambda *fiberadapter.FiberLambda

func init() {
	log.Printf("Fiber cold start")
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/create-workspace-body", func(c *fiber.Ctx) error {
		payload := struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		return c.JSON(payload)

	})

	app.Post("/create-workspace/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("👴 %s is %s years old", c.Query("name"), c.Query("age"))
		return c.SendString(msg) // => 👴 john is 75 years old
	})

	fiberLambda = fiberadapter.New(app)
}

// Handler will deal with Fiber working with Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return fiberLambda.ProxyWithContext(ctx, req)
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(Handler)
}
