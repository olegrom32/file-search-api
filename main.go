package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/olegrom32/file-search-api/internal"
	"github.com/olegrom32/file-search-api/internal/inputfile"
	loggerpkg "github.com/olegrom32/file-search-api/internal/logger"
	"github.com/olegrom32/file-search-api/internal/repository"
)

func main() {
	// TODO use a proper structured logger.
	// I didn't really understand this part of the task:
	// configuration file where you can specify log level (you should be able to choose between Info, Debug, Error)
	// Every log entry will have it's own log level. Maybe it asks to limit the amount of logs, like omit debug log
	// entries in production, but keep in dev, or something like that. But since I'm not using any fancy logger
	// (to keep it simple for this test task), there is no way to configure the std logger in this way.
	// So will use the configured logging entry prefix here.
	loggerPrefix := os.Getenv("LOGGER_PREFIX")
	if loggerPrefix == "" {
		loggerPrefix = "info"
	}

	logger, err := loggerpkg.New(loggerPrefix)
	if err != nil {
		panic(err)
	}

	// Create router
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger}))
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{}))

	// Init dependencies
	// Just reading config values directly from the env. Viper could be added here but I don't think it adds any value
	// in the context of this test task (our configuration is very simple).
	filename := os.Getenv("FILENAME")
	if filename == "" {
		filename = "input.txt"
	}

	margin := 0.1 // default 10%

	marginEnvStr := os.Getenv("MARGIN_PERCENT")
	if marginEnvStr != "" {
		marginEnv, err := strconv.Atoi(marginEnvStr)
		if err != nil {
			log.Fatal(err)
		}

		margin = float64(marginEnv) / 100
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	f, err := os.Open(filename)
	if err != nil {
		logger.Fatalf("Failed to open the input file %s: %s", filename, err.Error())
	}

	file, err := inputfile.Load(f)
	if err != nil {
		logger.Fatal(err)
	}

	fileRepo := repository.NewFileInMemory(file, margin)

	r.Get("/index/{value}", func(w http.ResponseWriter, r *http.Request) {
		// Always respond with json
		w.Header().Set("content-type", "application/json")

		valueStr := chi.URLParam(r, "value")

		value, err := strconv.Atoi(valueStr)
		if err != nil {
			logger.Printf("Invalid input parameter: %s", err.Error())

			writeError(w, http.StatusNotFound, err)

			return
		}

		res, err := fileRepo.FindByValue(value)
		if err != nil {
			logger.Printf("Failed to find the value index: %s", err.Error())

			switch {
			case errors.Is(err, internal.ErrNotFound):
				writeError(w, http.StatusNotFound, err)
			default:
				writeError(w, http.StatusInternalServerError, err)
			}

			return
		}

		// response is the response payload
		response := struct {
			Index int `json:"index"`
		}{
			Index: res,
		}

		// Marshal the response
		responseBytes, err := json.Marshal(response)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)

			return
		}

		// Send back the response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(responseBytes)
	})

	logger.Printf("Starting REST API on :%s", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Fatal(err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	if err == nil {
		return
	}

	w.WriteHeader(status)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"message":%q}`, err.Error())))
}
