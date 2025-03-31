package main

import (
	"GameApp/entity"
	"GameApp/repository/mysql"
	"fmt"
)

func main() {
	mysqlRepo := mysql.New()
	createdUser, err := mysqlRepo.Register(entity.User{
		ID:          0,
		PhoneNumber: "0912",
		Name:        "Hamed",
	})
	if err != nil {
		fmt.Errorf("cannot register user: %w", err)

	}
	fmt.Println(createdUser)

}
