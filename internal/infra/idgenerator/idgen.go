package idgenerator

import gonanoid "github.com/matoous/go-nanoid/v2"

type NanoIDGenerator struct{}

func (n NanoIDGenerator) Generate() string {
	id, err := gonanoid.New(7)
	if err != nil {
		panic(err)
	}

	return id
}
