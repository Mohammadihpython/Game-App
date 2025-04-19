package conf

var defaultConfig = map[string]interface{}{
	"auth.refresh_subject":    RefreshSubject,
	"auth.access_subject":     AccessSubject,
	"access_expiration_time":  AccessExpirationTime,
	"refresh_expiration_time": RefreshExpirationTime,
}
