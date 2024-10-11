package pkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

var ErrTokenExpired = errors.New("Token has expired")

type Maker interface {
	CreateToken(userID uint32, email string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

var _ Maker = (*PasetoMaker)(nil)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPaseto(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("SymmetricKey too short should be: %v", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(userID uint32, email string, duration time.Duration) (string, error) {
	payload, err := NewPayload(userID, email, duration)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
