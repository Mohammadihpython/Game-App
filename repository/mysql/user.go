package mysql

import (
	"GameApp/entity"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
)

func (d MYSQL) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	row := d.db.QueryRow("SELECT * FROM users WHERE phone_number=?", phoneNumber)
	_, err := scanRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, fmt.Errorf("cant scan query result: %w", err)
	}

	return false, nil
}

func (d MYSQL) RegisterUser(u entity.User) (entity.User, error) {
	res, err := d.db.Exec(`insert into users(Name,phone_number,password) values(?,?,?) `, u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("cant execute command : %w", err)

	}
	//	error is always nil
	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}
func (d MYSQL) GetUserByPhone(phone_number string) (entity.User, bool, error) {
	row := d.db.QueryRow(`SELECT * FROM users WHERE phone_number = ?`, phone_number)
	user, err := scanRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, false, nil
		}
		return entity.User{}, false, fmt.Errorf("cant scan query result: %w", err)
	}

	return user, true, nil
}

func (d MYSQL) GetUserByID(userid uint) (entity.User, error) {
	row := d.db.QueryRow("SELECT * FROM users WHERE id = ?", userid)
	user, err := scanRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, fmt.Errorf("cant scan query result: %w", err)

		}
		return entity.User{}, fmt.Errorf("cant scan query result: %w", err)
	}
	return user, nil
}
func scanRow(row *sql.Row) (entity.User, error) {
	user := entity.User{}
	var createdAT []uint8
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAT, &user.Password)
	return user, err

}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])

}
