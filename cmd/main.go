package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/isdzulqor/donation-hub/internal/driver/rest"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"strings"
)

type config struct {
	port         string
	dbDriverName string
	dbDataSource string
}

type customServeMux struct {
	*http.ServeMux
	excludedURLs []string
}

func main() {
	cfg := envConfig()
	db, _ := GetDatabaseConnection(cfg.dbDriverName, cfg.dbDataSource)

	//userDataStorage := userstorage.New(db)
	//userService := user.NewService(userDataStorage)

	api := rest.API{
		DB:             db,
		UserService:    nil,
		ProjectService: nil,
	}
	mux := &customServeMux{
		ServeMux: http.NewServeMux(),
		// exclude url for force write header Content-Type application/json
		excludedURLs: []string{
			"/",
			"favicon.ico",
			"/assets",
		},
	}
	mux.HandleFunc("/users/register", api.HandlePostUserRegister)
	mux.HandleFunc("/users/login", api.HandlePostUserLogin)
	mux.HandleFunc("/users", api.HandleGetUser)
	mux.HandleFunc("/projects/upload", api.HandleGetProjectUpload)
	mux.HandleFunc("/projects", api.HandleGetAndPostProject)
	mux.HandleFunc("/projects/", api.HandleGetProjectById)
	mux.HandleFunc("/projects/{project_id}/review", api.HandlePutProjectReview)
	mux.HandleFunc("/projects/{project_id}/donations", api.HandleGetAndPostProjectDonation)
	log.Println("Starting server on :" + cfg.port)
	err := http.ListenAndServe(":"+cfg.port, mux)
	log.Fatal(err)
}

func toJsonMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func (m *customServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestHandler, _ := m.ServeMux.Handler(r)
	found := false
	for _, url := range m.excludedURLs {
		if strings.HasPrefix(r.URL.Path, url) {
			found = true
			break
		}
	}

	// Apply middleware if not excluded and handler exists.
	if !found && requestHandler != nil {
		middleware := toJsonMiddleware(http.HandlerFunc(requestHandler.ServeHTTP))
		middleware(w, r)
	} else if !found {
		http.NotFound(w, r)
	} else {
		requestHandler.ServeHTTP(w, r)
	}
}

func envConfig() config {
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

	return config{port, dbDriverName, dbDataSource}
}

func GetDatabaseConnection(driverName string, dataSource string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driverName, dataSource)
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
