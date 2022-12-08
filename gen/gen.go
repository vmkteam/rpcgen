package gen

const version = "2.4.0"

const DefinitionsPrefix = "#/definitions/"

type GeneratorData struct {
	Version string
}

func DefaultGeneratorData() GeneratorData {
	return GeneratorData{
		Version: version,
	}
}
