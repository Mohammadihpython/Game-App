package authorizationservice

import (
	"GameApp/entity"
	"GameApp/pkg/richerror"
)

type Repository interface {
	GetUserPermissionTitles(userID uint, role entity.Role) ([]entity.PermissionTitle, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) CheckAccess(userID uint, role entity.Role, permissions ...entity.PermissionTitle) (bool, error) {
	// get the mysqluser role
	//get all ACLs for the given role
	// get all ACLs for the given mysqluser
	// merge all ACLs
	// the bow things Handel in repository

	//	 check the access
	const op = "authorizationservice.CheckAccess"
	permissionsTitles, err := s.repo.GetUserPermissionTitles(userID, role)
	if err != nil {
		return false, richerror.New(op).WithWrappedError(err).WithMessage("user not allowed")
	}

	for _, pT := range permissionsTitles {
		for _, p := range permissions {
			if p == pT {
				return true, nil
			}
		}
	}

	return false, nil

}
