package hasher

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type bcryptHasher struct {
	cost int
}

func New() Hasher {
	return &bcryptHasher{cost: 12}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (h *bcryptHasher) Compare(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}


