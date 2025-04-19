package conf

var defaultConfig = map[string]interface{}{
	"auth.refresh_subject":             RefreshSubject,
	"auth.access_subject":              AccessSubject,
	"auth.access_expiration_time":      AccessExpirationTime,
	"auth.refresh_expiration_time":     RefreshExpirationTime,
	"matching_service.waiting_timeout": WaitingTimeout,
}
