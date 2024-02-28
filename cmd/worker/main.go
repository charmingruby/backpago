package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/charmingruby/backpago/internal/bucket"
	"github.com/charmingruby/backpago/internal/queue"
)

func main() {
	// rabbitmq config
	qcfg := queue.RabbitMQConfig{
		URL:       "amqp://" + os.Getenv("RABBIT_URL"),
		TopicName: os.Getenv("RABBIT_TOPIC_NAME"),
		Timeout:   time.Second * 30,
	}

	// create new queue
	qc, err := queue.New(queue.RabbitMQ, qcfg)
	if err != nil {
		panic(err)
	}

	// create channel to consume messages
	c := make(chan queue.QueueDto, 1)
	go qc.Consume(c)

	// bucket config
	bcfg := bucket.AwsConfig{
		Config: &aws.Config{
			Region:      aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_KEY"), os.Getenv("AWS_SECRET"), ""),
		},
		BucketDownload: "aprenda-golang-drive-raw",
		BucketUpload:   "aprenda-golang-drive-gzip",
	}

	// create new bucket session
	b, err := bucket.New(bucket.AwsProvider, bcfg)
	if err != nil {
		panic(err)
	}

	log.Println("waiting for messages")
	for msg := range c {
		dst := fmt.Sprintf("%d_%s", msg.ID, msg.Filename)

		log.Printf("Start working on %s\n", msg.Filename)

		err := b.Download(msg.Path, dst)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		file, err := os.Open(dst)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		body, err := io.ReadAll(file)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		_, err = zw.Write(body)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		if err := zw.Close(); err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		zr, err := gzip.NewReader(&buf)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		err = b.Upload(zr, msg.Path)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		err = os.Remove(dst)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		log.Printf("%s was proccesed with success!\n", msg.Filename)
	}
}
