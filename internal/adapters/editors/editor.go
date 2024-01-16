package editors

type Editor interface {
	Open(string) error
}
