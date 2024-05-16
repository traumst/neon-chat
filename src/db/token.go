package db

type TokenType string

const (
	SignupToken TokenType = "signup"
)

type Token struct {
	Id    int       `db:"id"`
	Type  TokenType `db:"type"`
	Value string    `db:"value"`
}
