package core

type FinderInterface interface {
	Find(path string) (bool, error)
}
