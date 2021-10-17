package main

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {
	app := fiber.New()

	app.Post("/", processImage)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalln(err)
	}
}

func processImage(ctx *fiber.Ctx) error {
	var json = new(map[string]interface{})
	if err := ctx.BodyParser(json); err != nil {
		log.Fatalln(err)
		return fiber.ErrInternalServerError
	}

	imageB64 := (*json)["content"]
	if imageB64 == nil {
		return fiber.ErrBadRequest
	}

	file, err := os.CreateTemp("", "*.jpg")
	if err != nil {
		log.Fatalln(err)
		return fiber.ErrInternalServerError
	}

	image, err := base64.StdEncoding.DecodeString(imageB64.(string))
	if err != nil {
		log.Fatalln(err)
		return fiber.ErrUnprocessableEntity
	}

	if _, err := file.Write(image); err != nil {
		log.Fatalln(err)
		return fiber.ErrInternalServerError
	}
	if err := file.Close(); err != nil {
		log.Fatalln(err)
		return fiber.ErrInternalServerError
	}

	return detectTextOnImage(ctx, file.Name())
}

func detectTextOnImage(c *fiber.Ctx, filename string) error {
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalln(err)
		return fiber.ErrInternalServerError
	}
	defer client.Close()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
		return fiber.ErrInternalServerError
	}
	defer file.Close()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalln(err)
		return fiber.ErrInternalServerError
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("DetectTexts: %v\n", err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"annotations": annotations,
	})
}
