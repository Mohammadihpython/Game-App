package conf

import "time"

const (
	SECRET                = "Hmdsfksdf"
	AccessExpirationTime  = time.Hour * 24
	RefreshExpirationTime = time.Hour * 24 * 7
	AccessSubject         = "at"
	RefreshSubject        = "rt"
)
