package shares

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type BaseEntity struct {
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	BaseEntity
	name string
}

func NewUserByName(name string) User {
	return User{
		name: name,
	}
}

func NewUserById(id int64) User {
	return User{
		BaseEntity: BaseEntity{
			Id: id,
		},
	}
}

func Insert[T BaseEntity](t T) {
}

type Stream[T any] struct {
}

func (s *Stream[T]) Filter() {
}

type Selector[T any] struct {
}

func (s *Selector[T]) Get() (*T, error) {
	return new(T), nil
}

func TestUseSelector(t *testing.T) {
	s := &Selector[User]{}
	user, err := s.Get()
	assert.NoError(t, err)
	t.Log(user)
}
