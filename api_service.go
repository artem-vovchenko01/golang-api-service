package main

import (
	"log"
	"net/http"
	"context"
	"os"
	"os/signal"
	"time"
	"github.com/gorilla/mux"
	"strings"
)

func getCakeHandler(w http.ResponseWriter, r *http.Request, u User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(u.FavoriteCake))
}

func main() {
	r := mux.NewRouter()

	users := NewInMemoryUserStorage()
	userService := UserService {
		repository: users,
	}

	jwtService, err := NewJWTService("privkey.rsa", "pubkey.rsa")
	if err != nil {
		panic(err)
	}

	r.HandleFunc("/cake", logRequest(jwtService.jwtAuth(users, getCakeHandler))).Methods(http.MethodGet)
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

type ProtectedHandler func(rw http.ResponseWriter, r *http.Request, u User) 

func (j *JWTService) jwtAuth(
	users UserRepository,
	h ProtectedHandler,
) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		auth, err := j.ParseJWT(token)
		if err != nil {
			rw.WriteHeader(401)
			rw.Write([]byte("unathorized"))
			return
		}

		user, err := users.Get(auth.Email)
		if err != nil {
			rw.WriteHeader(401)
			rw.Write([]byte("unathorized"))
			return
		}

		h(rw, r, user)
	}
}
