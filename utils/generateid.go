package utils

import (
	"math/rand"
	"time"
)

type IdGenerator interface {
	GenerateId() (int64, error)
}

type BasicGenerator struct {
	rand rand.Rand
}

func NewBasicGenerator() *BasicGenerator {
	return &BasicGenerator{
		rand: *rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (g *BasicGenerator) GenerateId() (int64, error) {
	id := g.rand.Int63()
	return id, nil
}
