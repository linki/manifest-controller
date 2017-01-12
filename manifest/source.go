package manifest

type Source interface {
	Fetch() (changed bool, err error)
	Manifest() string
}
