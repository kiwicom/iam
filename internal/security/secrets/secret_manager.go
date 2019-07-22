package secrets

// SecretManager is a interface that describes how we want to use secrets
type SecretManager interface {
	GetAppToken(string, string) (string, error)
	GetSetting(string) (string, error)
}
