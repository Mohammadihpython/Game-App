package mysql

import (
	"GameApp/entity"
	"database/sql"
	"fmt"
)

func (d MYSQL) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	user := entity.User{}
	var createdAT []uint8
	row := d.db.QueryRow("SELECT * FROM users WHERE phonenumber=?", phoneNumber)
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAT)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, fmt.Errorf("cant scan query result: %w", err)
	}

	return false, nil
}

func (d MYSQL) RegisterUser(u entity.User) (entity.User, error) {
	res, err := d.db.Exec(`insert into usSers() values(?,?) `, u.Name, u.PhoneNumber)
	if err != nil {
		return entity.User{}, fmt.Errorf("cant execute command : %w", err)

	}
	//	error is always nil
	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}
