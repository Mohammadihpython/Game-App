package accesscontrol

import (
	"GameApp/entity"
	"GameApp/repository/mysql"
)

func scanPermission(scanner mysql.Scanner) (entity.Permission, error) {
	permission := entity.Permission{}
	var createdAT []uint8
	err := scanner.Scan(&permission.ID, &permission.Title, &createdAT)
	return permission, err

}
