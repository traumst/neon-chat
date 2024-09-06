package interfaces

// implements Identifiable interface for convenience
type Renderable interface {
	// returns raw html string representation of underlying object
	HTML() (string, error)
}
