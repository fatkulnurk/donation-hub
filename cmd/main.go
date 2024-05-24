package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/isdzulqor/donation-hub/internal/core/model"
	"github.com/isdzulqor/donation-hub/internal/core/service/project"
	"github.com/isdzulqor/donation-hub/internal/core/service/user"
	"github.com/isdzulqor/donation-hub/internal/driven/auth/jwt"
	"github.com/isdzulqor/donation-hub/internal/driven/storage/mysql/projectstorage"
	"github.com/isdzulqor/donation-hub/internal/driven/storage/mysql/userstorage"
	"github.com/isdzulqor/donation-hub/internal/driven/storage/s3/projectfilestorage"
	"github.com/isdzulqor/donation-hub/internal/driver/rest"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

func main() {
	cfg := envConfig()
	db, err := GetDatabaseConnection(cfg.DBDriverName, cfg.DBDataSource)
	if err != nil {
		log.Fatalln(err.Error())
	}

	s3Client, err := initializeS3Client()
	if err != nil {
		log.Fatalln(err.Error())
	}

	authToken := jwt.New(cfg)

	userStorage := userstorage.New(db)
	userService := user.NewService(userStorage, authToken)

	projectFileStorage := projectfilestorage.NewStorage(s3Client)
	projectStorage := projectstorage.New(db)
	projectService := project.NewService(projectStorage, projectFileStorage, userStorage)

	api := rest.API{
		DB:             db,
		UserService:    userService,
		ProjectService: projectService,
		AuthToken:      authToken,
	}

	api.ListenAndServe(&cfg)
}

func envConfig() model.ConfigMap {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		panic("APP_PORT not provided")
	}

	dbDriverName, ok := os.LookupEnv("DATABASE_DRIVER_NAME")
	if !ok {
		panic("DATABASE_DRIVER_NAME not provided")
	}

	dbDataSource, ok := os.LookupEnv("DATABASE_DATA_SOURCE")
	if !ok {
		panic("DATABASE_DATA_SOURCE not provided")
	}

	secretKey := "supersecrethehehe"
	issuer := "Donation Hub"

	return model.ConfigMap{
		Port:         port,
		DBDriverName: dbDriverName,
		DBDataSource: dbDataSource,
		SecretKey:    secretKey,
		Issuer:       issuer,
	}
}

func GetDatabaseConnection(driverName string, dataSource string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", dataSource)
	log.Println("Get Database Connection")
	log.Println(driverName)
	log.Println(dataSource)
	if err != nil {
		log.Fatalln(err)

		return nil, err
	}

	fmt.Println("Database Connected")

	return db, nil
}

func initializeS3Client() (s3Client *s3.Client, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_DEFAULT_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			os.Getenv("AWS_SESSION_TOKEN"),
		)),
	)

	if err != nil {
		return s3Client, fmt.Errorf("failed to load configuration, %w", err)
	}

	s3Client = s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(os.Getenv("LOCALSTACK_ENDPOINT"))
		options.UsePathStyle = os.Getenv("AWS_USE_PATH_STYLE_ENDPOINT") == "1"
	})

	return s3Client, err
}
