package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello brothers!"))
}

var data map[string]string = map[string]string{}

// Data will be of the follwing format:
// {
// "temp1" : 69,
// "temp1" : 69,
// "temp1" : 69,
// "temp1" : 69,
// "temp1" : 69,
// "temp1" : 69,
// "env_temp" : 69,
// "env_humid": 70,
// "env_lux": 34,
// "sht_humid": 5,
// "sht_temp": 5
// }
type DataModel struct {
	Temp1    float32 `json:"temp1"`
	Temp2    float32 `json:"temp2"`
	Temp3    float32 `json:"temp3"`
	Temp4    float32 `json:"temp4"`
	Temp5    float32 `json:"temp5"`
	Temp6    float32 `json:"temp6"`
	EnvTemp  float32 `json:"env_temp"`
	EnvHumid float32 `json:"env_humid"`
	EnvLux   float32 `json:"env_lux"`
	ShtHumid float32 `json:"sht_humid"`
	ShtTemp  float32 `json:"sht_temp"`
}

func WriteData(w http.ResponseWriter, r *http.Request) {
	// Declare a new Person struct.
	var d DataModel

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Do something with the Person struct...
	fmt.Printf("Data: %+v", d)
	// current timestamp
	curr := fmt.Sprintf("%d", time.Now().Unix())
	data[curr] = fmt.Sprintf("%+v", d)

	fmt.Printf("Data: %+v", data)

	w.WriteHeader(http.StatusOK)
	return
}

func ReadAllData(w http.ResponseWriter, r *http.Request) {
	// read data from DB

	d := data

	jsonEncoded, err := json.Marshal(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Print(string(jsonEncoded))

	w.Write(jsonEncoded)
	return
}

func (app *application) routes() http.Handler {
	// init the router
	router := chi.NewRouter()

	// TODO
	// router.NotFound()

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*", "*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	app.infoLog.Println("CORS enabled!")

	// for testing purposes, does not require JWT authorization
	// should not be used in production
	router.Route("/", func(r chi.Router) {
		r.Get("/", hello)
		r.Post("/data", WriteData)
		r.Get("/data", ReadAllData)

		// r.Route("/zip", func(r chi.Router) {

		// 	r.Get("/coordinates", app.getZipCodeCoords)
		// })
		// r.Route("/plot", func(r chi.Router) {
		// 	r.Get("/{plotId}", app.GetPlotById)
		// 	r.Put("/sync", app.SyncBumbalZones)
		// 	r.Post("/save", app.SavePlot)
		// 	r.Delete("/{plotId}", app.DeletePlotById)
		// })

	})

	// authorized routes
	// router.Group(func(r chi.Router) {
	// 	r.Use(JwtChecker)
	// 	r.Use(JWTRequestChecker)

	// 	r.Route("/zip", func(r chi.Router) {
	// 		r.Get("/coordinates", app.getZipCodeCoords)
	// 	})

	// 	r.Route("/user", func(r chi.Router) {
	// 		r.Get("/plots", app.getUserPlotIDs)
	// 	})

	// 	r.Route("/plot", func(r chi.Router) {
	// 		r.Get("/{plotId}", app.GetPlotById)
	// 		r.Put("/sync", app.SyncBumbalZones)
	// 		r.Post("/save", app.SavePlot)
	// 		r.Delete("/{plotId}", app.DeletePlotById)
	// 	})

	// 	r.Route("/bumbal", func(r chi.Router) {
	// 		r.Put("/algorithm/kmeans", bumbal.RunKMeans)
	// 		r.Put("/algorithm/genetic", bumbal.RunGenetic)
	// 	})

	// })

	// log all the routes mounted on the router
	app.infoLog.Println("Mounted routes:")
	err := chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		app.infoLog.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})
	if err != nil {
		panic(err)
	}

	return router
}

func main() {
	addr := ":3000"

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)

	srv.ListenAndServe()
}
