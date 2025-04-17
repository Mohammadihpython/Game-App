package entity

// AccessControl only keeps allowed permission
type AccessControl struct {
	ID uint
	// I used database polymerize
	ActorID      uint
	ActorType    ActorType
	PermissionID uint
}

type ActorType string

const (
	RoleActorType = "role"
	UserActorType = "mysqluser"
)
