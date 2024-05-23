package asset

type IReaderLock interface {
	Setup() error
	Reload() error
	Get(name string) ([]byte, error)
	GetHash(name string) ([]byte, string, error)
	Names() []string
}
