package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/charmingruby/backpago/internal/auth"
	"github.com/charmingruby/backpago/internal/bucket"
	"github.com/charmingruby/backpago/internal/files"
	"github.com/charmingruby/backpago/internal/folders"
	"github.com/charmingruby/backpago/internal/queue"
	"github.com/charmingruby/backpago/internal/users"
	"github.com/charmingruby/backpago/pkg/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	db, b, qc := getSessions()

	r := chi.NewRouter()

	// allow cors
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Accept", "Content-Type"},
	}))

	// define endpoints
	r.Post("/auth", auth.HandleAuth(func(login, pass string) (auth.Authenticated, error) {
		return users.Authenticate(login, pass)
	}))

	files.SetRoutes(r, db, b, qc)
	folders.SetRoutes(r, db)
	users.SetRoutes(r, db)

	// start server
	fmt.Println("Server start")
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")), r)
}

func getSessions() (*sql.DB, *bucket.Bucket, *queue.Queue) {
	// create new database connection
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal(err)
	}

	// rabbitmq config
	qcfg := queue.RabbitMQConfig{
		URL:       "amqp://" + os.Getenv("RABBIT_URL"),
		TopicName: os.Getenv("RABBIT_TOPIC_NAME"),
		Timeout:   time.Second * 30,
	}

	// create new queue
	qc, err := queue.New(queue.RabbitMQ, qcfg)
	if err != nil {
		log.Fatal(err)
	}

	// bucket config
	bcfg := bucket.AwsConfig{
		Config: &aws.Config{
			Region:      aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_KEY"), os.Getenv("AWS_SECRET"), ""),
		},
		BucketDownload: "aprenda-golang-drive-gzip",
		BucketUpload:   "aprenda-golang-drive-raw",
	}

	// create new bucket session
	b, err := bucket.New(bucket.AwsProvider, bcfg)
	if err != nil {
		log.Fatal(err)
	}

	return db, b, qc
}
