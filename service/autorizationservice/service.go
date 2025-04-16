package autorizationservice

type Service struct {
}

func (s Service) CheckAccess(userID uint) (bool, error) {
	// get the user role
	//get all ACLs for the given role
	// get all ACLs for the given user
	// merge all ACLs

	//	 check the access

}
