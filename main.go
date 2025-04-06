package main

import (
	"GameApp/repository/mysql"
	"GameApp/service/authservice"
	"GameApp/service/userservice"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const SECRET = "Hmdsfksdf"

func main() {

	//http.HandleFunc("/users/register",userRegisterHandler)
	//http.ListenAndServe(:8080,nil)
	// another way is to use multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/users/register", userRegisterHandler)
	mux.HandleFunc("/health", HealthCheck)
	mux.HandleFunc("/users/login", userLoginHandler)
	mux.HandleFunc("/users/profile", userProfileHandler)

	server := http.Server{Addr: ":8090", Handler: mux}
	log.Fatal(server.ListenAndServe())

}
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi")
	fmt.Fprintf(w, "healthy ")

}
func userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, `{"error":"invalid Method" }`)
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}
	var req userservice.RegisterRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))

		return
	}
	mysqlRepo := mysql.New()
	authSvc := authservice.New(SECRET, time.Hour*24, time.Hour*24*7, "at", "rt")

	userSvc := userservice.New(mysqlRepo, authSvc)
	_, err = userSvc.Register(req)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`{"error:" "%s"}`, err.Error())))

		return
	}
	w.Write([]byte(`"{"message":"user created"}"`))

}
func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, `{"error":"invalid Method" }`)
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
	}
	var req userservice.LoginRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}
	mysqlRepo := mysql.New()
	authSvc := authservice.New(SECRET, time.Hour*24, time.Hour*24*7, "at", "rt")

	userSvc := userservice.New(mysqlRepo, authSvc)
	res, err := userSvc.Login(req)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}
	data, err = json.Marshal(res)
	fmt.Println(string(data))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))

		return
	}
	w.Write(data)

}
func userProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, `{"error":"invalid Method" }`)

	}
	authSvc := authservice.New(SECRET, time.Hour*24, time.Hour*24*7, "at", "rt")
	// Get user id from Jwt Token
	auth := r.Header.Get("Authorization")
	claims, err := authSvc.ParsToken(auth)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	pReq := userservice.ProfileRequest{UserID: claims.UserID}
	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo, authSvc)
	res, err := userSvc.Profile(pReq)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}
	fmt.Println(res)
	data, err := json.Marshal(res)
	fmt.Println(string(data))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))

		return
	}
	w.Write(data)

}
