package main

import (
	"log"
	"net/http"
	"context"
	"os"
	"os/signal"
	"time"
	"github.com/gorilla/mux"
)

func getCakeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("cake"))
}

func main() {
	r := mux.NewRouter()

	userService := UserService {
		repository: NewInMemoryUserStorage(),
	}


	jwtService, err := NewJWTService("privkey.rsa", "pubkey.rsa")
	if err != nil {
		panic(err)
	}

	r.HandleFunc("/cake", logRequest(getCakeHandler)).Methods(http.MethodGet)
	r.HandleFunc("/user/register", logRequest(userService.Register)).Methods(http.MethodPost)
	r.HandleFunc("/user/jwt", logRequest(wrapJWT(jwtService, userService.JWT))).Methods(http.MethodPost)
	srv := http.Server {
		Addr: 		":8080",
		Handler: 	r,
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<- interrupt
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	log.Println("Server started, hbit Ctrl+C to stop")
	err = srv.ListenAndServe()
	if err != nil {
		log.Println("Server exited with error", err)
	}

	log.Println("Good bye :)")
}

func wrapJWT(
	jwt *JWTService,
	f func(http.ResponseWriter, *http.Request, *JWTService),
) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		f(rw, r, jwt)
	}
}

