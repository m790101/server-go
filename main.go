package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"server/internal/database"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	Db             *database.DB
	Secret         string
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func main() {

	godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")

	// fmt.Println(ss, err)

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	port := "8080"
	mux := http.NewServeMux()
	filepathRoot := "."
	apiCfg := apiConfig{fileserverHits: 0, Db: db, Secret: jwtSecret}
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerReset)
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleGetOne)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirps)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handleUpdateUser)
	mux.HandleFunc("POST /api/login", apiCfg.Login)
	mux.HandleFunc("GET /api/users", apiCfg.handleGetUsers)
	mux.HandleFunc("GET /healthz", handleHealthz)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerValidate)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}

func handleHealthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}
