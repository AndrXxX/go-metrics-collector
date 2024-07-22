package hashgenerator

type hashGenerator struct {
	key string
}

func New(key string) *hashGenerator {
	return &hashGenerator{key: key}
}

func (h *hashGenerator) Generate(data []byte) (string, error) {
	// TODO:
	return string(data), nil
}
