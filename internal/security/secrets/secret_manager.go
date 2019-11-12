package secrets

// SecretManager is a interface that describes how we want to use secrets
type SecretManager interface {
	DoesTokenExist(string) bool
	GetSetting(string) (string, error)
}
