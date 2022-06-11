package disorder_test

import (
	"fmt"
	"testing"

	"github.com/meerkat-lib/disorder"
	"github.com/meerkat-lib/disorder/internal/generator/golang"
	"github.com/meerkat-lib/disorder/internal/loader"
)

func TestLoadYamlFile(t *testing.T) {
	loader := loader.NewYamlLoader()
	files, err := loader.Load("./internal/test_data/schema.yaml")
	fmt.Println(err)
	generator := golang.NewGoGenerator()
	err = generator.Generate("./internal", files)
	fmt.Println(err)
	t.Fail()
}

func TestLoadJsonFile(t *testing.T) {
	loader := loader.NewJsonLoader()
	files, err := loader.Load("./internal/test_data/schema.json")
	fmt.Println(err)
	generator := golang.NewGoGenerator()
	err = generator.Generate("./internal", files)
	fmt.Println(err)
	t.Fail()
}

type S struct {
	Value uint8
}

func TestMarshal(t *testing.T) {

	input := map[string]string{}
	/*
		input := S{
			Value: 123,
		}*/
	data, err := disorder.Marshal(input)
	fmt.Println(err)
	fmt.Println(input)
	fmt.Println(data)
	var output interface{}
	err = disorder.Unmarshal(data, &output)
	fmt.Println(err)
	fmt.Println(output)

	t.Fail()
}
