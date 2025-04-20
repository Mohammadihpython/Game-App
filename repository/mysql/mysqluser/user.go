package mysqluser

import (
	"GameApp/entity"
	"GameApp/pkg/richerror"
	"GameApp/repository/mysql"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
)

func (d *DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	row := d.conn.Conn().QueryRow("SELECT * FROM users WHERE phone_number=?", phoneNumber)
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

func (d *DB) RegisterUser(u entity.User) (entity.User, error) {
	res, err := d.conn.Conn().Exec(`insert into users(Name,phone_number,password,role) values(?,?,?,?) `, u.Name, u.PhoneNumber, u.Password, u.Role.String())
	if err != nil {
		return entity.User{}, fmt.Errorf("cant execute command : %w", err)

	}
	//	error is always nil
	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}
func (d *DB) GetUserByPhone(phone_number string) (entity.User, error) {
	row := d.conn.Conn().QueryRow(`SELECT * FROM users WHERE phone_number = ?`, phone_number)
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

func (d *DB) GetUserByID(ctx context.Context, userid uint) (entity.User, error) {

	// we use QueryRowContext to know if the context is close we do not query
	row := d.conn.Conn().QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", userid)
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
	var rolstr string
	err := scanner.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAT, &user.Password, &rolstr)
	user.Role = entity.MapToRoleEntity(rolstr)
	return user, err

}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])

}
