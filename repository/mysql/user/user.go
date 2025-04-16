package user

import (
	"GameApp/entity"
	"GameApp/pkg/richerror"
	"GameApp/repository/mysql"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
)

func (d *mysql.MYSQL) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	row := d.db.QueryRow("SELECT * FROM users WHERE phone_number=?", phoneNumber)
	_, err := scanRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, richerror.New("mysql.IsPhoneNumberUnique").WithMessage("can scan query result").
			WithKind(richerror.KindUnexpected)
	}

	return false, nil
}

func (d *mysql.MYSQL) RegisterUser(u entity.User) (entity.User, error) {
	res, err := d.db.Exec(`insert into users(Name,phone_number,password) values(?,?,?) `, u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("cant execute command : %w", err)

	}
	//	error is always nil
	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}
func (d *mysql.MYSQL) GetUserByPhone(phone_number string) (entity.User, error) {
	row := d.db.QueryRow(`SELECT * FROM users WHERE phone_number = ?`, phone_number)
	user, err := scanRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, richerror.New("userservice.GetUserByPhone").WithWrappedError(err).
				WithKind(richerror.KindNotFound)
		}
		return entity.User{}, richerror.New("userservice.GetUserByPhone").WithWrappedError(err).
			WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d *mysql.MYSQL) GetUserByID(userid uint) (entity.User, error) {
	row := d.db.QueryRow("SELECT * FROM users WHERE id = ?", userid)
	user, err := scanRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, richerror.New("userservice.GetUserByID").WithWrappedError(err)

		}
		return entity.User{}, richerror.New("userservice.GetUserByID").WithWrappedError(err)
	}
	return user, nil
}
func scanRow(scanner mysql.Scanner) (entity.User, error) {
	// ParseTime=true handel fileds that time.time type and we didnt meed to convert to
	// []byte like var createdAT []uint8 instead we use time.time
	user := entity.User{}
	var createdAT []uint8
	err := scanner.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAT, &user.Password)
	return user, err

}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])

}
