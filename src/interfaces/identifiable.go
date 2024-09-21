package interfaces

// an identifaible entity, tycaly sourced from the database
type Identifiable interface {
	// db issued id, unique per type of entity
	GetId() uint
}
