package accesscontrol

import (
	"GameApp/entity"
	"GameApp/pkg/richerror"
	"GameApp/pkg/slice"
	"GameApp/repository/mysql"
	"strings"
)

func (d *db) GetUserPermissionTitles(userID uint) ([]entity.PermissionTitle, error) {
	const op = "mysql.GetUserPermissionTitles"
	//	get user
	user, err := d.GetUserByID(userID)
	if err != nil {
		return nil, richerror.New(op).WithWrappedError(err)

	}
	roleACL := make([]entity.AccessControl, 0)

	rows, err := d.db.Query("SELECT * FROM access_control WHERE actor_type= 'role' and actor_id=?", user.Role)
	if err != nil {
		return nil, richerror.New(op).WithWrappedError(err).
			WithKind(richerror.KindUnexpected).
			WithMessage("an unexpected error occurred")
	}

	defer rows.Close()

	for rows.Next() {
		acl, err := scanAccessControl(rows)
		if err != nil {
			return nil, richerror.New(op).WithWrappedError(err)

		}
		roleACL = append(roleACL, acl)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New(op).WithWrappedError(err)
	}

	userACL := make([]entity.AccessControl, 0)
	userRows, err := d.db.Query("SELECT * FROM access_control WHERE actor_type= 'user' and actor_id=?", user.ID)
	if err != nil {
		return nil, richerror.New(op).WithWrappedError(err).
			WithKind(richerror.KindUnexpected).
			WithMessage("an unexpected error occurred")
	}

	defer userRows.Close()

	for userRows.Next() {
		acl, err := scanAccessControl(userRows)
		if err != nil {
			return nil, richerror.New(op).WithWrappedError(err)

		}
		userACL = append(userACL, acl)
	}
	defer userRows.Close()

	if err := userRows.Err(); err != nil {
		return nil, richerror.New(op).WithWrappedError(err)
	}
	//	 merge acls by permission id
	permissionIDs := make([]uint, 0)
	for _, r := range roleACL {
		if !slice.DoseExist(permissionIDs, r.PermissionID) {
			permissionIDs = append(permissionIDs, r.PermissionID)
		}

	}
	if len(permissionIDs) == 0 {
		return nil, nil
	}
	//	select * from permission where id in(?,?,?)
	//	select * from permission where id in(1,2,5)
	args := make([]interface{}, len(permissionIDs))
	for i, id := range permissionIDs {
		args[i] = id
	}
	rows, err = d.db.Query(
		"select * from permission where id in (?"+strings.Repeat(",?", len(permissionIDs)-1)+")", args...)
	if err != nil {
		return nil, richerror.New(op).WithWrappedError(err)
	}
	defer rows.Close()

	permissionTitles := make([]entity.PermissionTitle, 0)
	for rows.Next() {
		permission, err := scanPermission(rows)
		if err != nil {
			return nil, richerror.New(op).WithWrappedError(err)
		}
		permissionTitles = append(permissionTitles, permission.Title)

	}
	if err := rows.Err(); err != nil {
		return nil, richerror.New(op).WithWrappedError(err)
	}
	return permissionTitles, nil

}

func scanAccessControl(scanner mysql.Scanner) (entity.AccessControl, error) {
	// ParseTime=true handel fileds that time.time type and we didnt meed to convert to
	// []byte like var createdAT []uint8 instead we use time.time
	acl := entity.AccessControl{}
	var createdAT []uint8
	err := scanner.Scan(&acl.ID, &acl.ActorID, &acl.ActorType, acl.PermissionID, &createdAT)
	return acl, err

}
