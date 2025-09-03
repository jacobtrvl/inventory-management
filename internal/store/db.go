package store

// Store interface.
// This is not a idiomatic Go style to define interface when there is only one implementation.
// This is just for demonstrating usage of interfaces.
type Store interface {
	Write(table string, key any, item any) error
	Read(table string, id any) (any, error)
	ReadRange(table string, start, end int) ([]any, bool, error)
	ReadAll(table string) ([]any, error)
	Delete(table string, id any) error
	CreateTable(name string) error
	DeleteTable(name string) error
}
