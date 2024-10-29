package hashgenerator

type hashGeneratorFactory struct {
}

// Factory возвращает фабрику для получения сервисов генерации хешей
func Factory() *hashGeneratorFactory {
	return &hashGeneratorFactory{}
}

// SHA256 возвращает сервисов генерации хешей на основе алгоритма SHA256
func (f *hashGeneratorFactory) SHA256() *sha256Generator {
	return &sha256Generator{}
}
