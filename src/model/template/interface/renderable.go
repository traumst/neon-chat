package interfaces

type Renderable interface {
	HTML() (string, error)
}
