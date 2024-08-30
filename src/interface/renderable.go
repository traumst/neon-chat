package interfaces

// implements Identifiable interface for convenience
type Renderable interface {
	// returns raw html string representation of underlying object
	HTML() (string, error)
	// returns briferer html string representation of underlying object
	ShortHTML() (string, error)
}
