package autorizationservice

import "GameApp/entity"

type Repository interface {
	GetUserPermissionTitles(userID uint) ([]entity.PermissionTitle, error)
}

type Service struct {
}

func (s Service) CheckAccess(userID uint, permissions ...entity.PermissionTitle) (bool, error) {
	// get the user role
	//get all ACLs for the given role
	// get all ACLs for the given user
	// merge all ACLs

	//	 check the access

	return false, nil

}
