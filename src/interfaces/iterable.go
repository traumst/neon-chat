package interfaces

type Iterable interface {
	Next() (bool, Identifiable)
}
