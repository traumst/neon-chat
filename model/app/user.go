package app

type User struct {
	Id   uint
	Name string
	Salt []byte
}
