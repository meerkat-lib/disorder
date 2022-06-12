package disorder_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/meerkat-lib/disorder"
	"github.com/meerkat-lib/disorder/internal/generator/golang"
	"github.com/meerkat-lib/disorder/internal/loader"
	"github.com/meerkat-lib/disorder/internal/test_data/test/sub"
	"github.com/meerkat-lib/disorder/rpc"
)

type mathService struct {
}

func (*mathService) Increase(c *rpc.Context, request *int32) (*int32, *rpc.Error) {
	value := *request
	value++
	return &value, nil
}

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

func TestLoadTomlFile(t *testing.T) {
	loader := loader.NewTomlLoader()
	files, err := loader.Load("./internal/test_data/schema.toml")
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

func TestRpc(t *testing.T) {
	s := rpc.NewServer()
	sub.RegisterMathService(s, &mathService{})
	err := s.Listen(":8080")
	fmt.Println(err)

	time.Sleep(time.Second)

	c := sub.NewMathServiceClient(rpc.Dial("localhost:8080"))
	value := int32(1)
	result, rpcErr := c.Increase(rpc.NewContext(), &value)
	fmt.Println(rpcErr)
	fmt.Println(*result)
	t.Fail()
}
