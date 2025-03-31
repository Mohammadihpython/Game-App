package main

import (
	"GameApp/entity"
	"GameApp/repository/mysql"
	"GameApp/service/userservice"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {

	//http.HandleFunc("/users/register",userRegisterHandler)
	//http.ListenAndServe(:8080,nil)
	// another way is to use multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/users/register", userRegisterHandler)
	mux.HandleFunc("/health", HealthCheck)
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
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}
	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo)
	_, err = userSvc.Register(req)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error:" "%s"}`, err.Error())))
	}
	w.Write([]byte(`"{"message":"user created"}"`))

}
func Testmysqlconector() {
	mysqlRepo := mysql.New()
	createdUser, err := mysqlRepo.RegisterUser(entity.User{
		ID:          0,
		PhoneNumber: "0912",
		Name:        "Hamed",
	})
	if err != nil {
		fmt.Errorf("cannot register user: %w", err)
	}
	fmt.Println(createdUser)
}
