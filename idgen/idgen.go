package idgen

import (
	"github.com/rs/xid"
	"github.com/satori/go.uuid"
)

type generator func() string

func xidGen() string {
	return xid.New().String()
}

func uuidGen() string {
	return uuid.NewV4().String()
}

func NewIdGenerator() *IdGenerator {
	return &IdGenerator{generator: xidGen}
}

func NewUUIDGenerator() *IdGenerator {
	return &IdGenerator{generator: uuidGen}
}

type IdGenerator struct {
	generator generator
}

func (z *IdGenerator) NewID() string {
	return z.generator()
}
