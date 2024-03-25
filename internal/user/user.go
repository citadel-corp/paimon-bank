package user

type User struct {
	ID             uint64
	Email          string
	Name           string
	HashedPassword string
}
