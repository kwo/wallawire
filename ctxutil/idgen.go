package ctxutil

import (
	"github.com/kwo/uuidv6"
)

func NewIdGenerator() (*IdGenerator, error) {

	generator, errGenerator := uuidv6.NewGenerator()
	if errGenerator != nil {
		return nil, errGenerator
	}

	return &IdGenerator{
		generator: generator,
	}, nil

}

type IdGenerator struct {
	generator *uuidv6.Generator
}

func (z *IdGenerator) NewID() string {
	return z.generator.New().String()
}
