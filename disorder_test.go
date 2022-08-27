package disorder_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/meerkat-io/disorder"
	"github.com/meerkat-io/disorder/internal/generator/golang"
	"github.com/meerkat-io/disorder/internal/loader"
	"github.com/meerkat-io/disorder/internal/test_data/test"
	"github.com/meerkat-io/disorder/internal/test_data/test/sub"
	"github.com/meerkat-io/disorder/rpc"
	"github.com/stretchr/testify/assert"
)

type testService struct {
}

func (*testService) Increase(c *rpc.Context, request int32) (int32, *rpc.Error) {
	fmt.Printf("input value: %d\n", request)
	request++
	return request, nil
}

func (*testService) PrintObject(c *rpc.Context, request *test.Object) (*test.Object, *rpc.Error) {
	fmt.Printf("input object: %v\n", *request)
	t := time.Now()
	color := test.ColorGreen
	return &test.Object{
		Time:        &t,
		IntField:    456,
		StringField: "bar",
		EnumField:   &color,
		IntArray:    []int32{4, 5, 6},
		IntMap: map[string]int32{
			"4": 4,
			"5": 5,
			"6": 6,
		},
		ObjArray: []*sub.SubObject{{Value: 456}},
		ObjMap: map[string]*sub.SubObject{
			"bar": {Value: 456},
		},
	}, nil
}

func (*testService) PrintSubObject(c *rpc.Context, request *sub.SubObject) (*sub.SubObject, *rpc.Error) {
	fmt.Printf("input sub object: %v\n", *request)
	return &sub.SubObject{
		Value: 456,
	}, nil
}

func (*testService) PrintTime(c *rpc.Context, request *time.Time) (*time.Time, *rpc.Error) {
	fmt.Printf("input time: %v\n", *request)
	t := time.Now()
	return &t, nil
}

func (*testService) PrintArray(c *rpc.Context, request []int32) ([]int32, *rpc.Error) {
	fmt.Printf("input array: %v\n", request)
	return []int32{4, 5, 6}, nil
}

func (*testService) PrintEnum(c *rpc.Context, request *test.Color) (*test.Color, *rpc.Error) {
	reqColor, _ := request.ToString()
	fmt.Printf("input enum: %s\n", reqColor)
	color := test.ColorGreen
	return &color, nil
}

func (*testService) PrintMap(c *rpc.Context, request map[string]string) (map[string]string, *rpc.Error) {
	fmt.Printf("input map: %s\n", request)
	return map[string]string{
		"bar": "foo",
	}, nil
}

func TestLoadYamlFile(t *testing.T) {
	loader := loader.NewYamlLoader()
	files, err := loader.Load("./internal/test_data/schema.yaml")
	assert.Nil(t, err)

	generator := golang.NewGoGenerator()
	err = generator.Generate("./internal", files)
	assert.Nil(t, err)
}

func TestMarshal(t *testing.T) {
	tt := time.Now()
	color := test.ColorRed
	target := test.Object{
		Time:        &tt,
		IntField:    123,
		StringField: "foo",
		EnumField:   &color,
		IntArray:    []int32{1, 2, 3},
		IntMap: map[string]int32{
			"1": 1,
			"2": 2,
			"3": 3,
		},
		ObjArray: []*sub.SubObject{{Value: 123}},
		ObjMap: map[string]*sub.SubObject{
			"foo": {Value: 123},
		},
		AnyField: "some text",
		AnyArray: []interface{}{"abc", 123, 3.14},
		AnyMap:   map[string]interface{}{"a": "abc", "b": "123", "c": 3.14},
	}
	data, err := disorder.Marshal(&target)
	fmt.Println(data)
	assert.Nil(t, err)

	var result1 interface{}
	err = disorder.Unmarshal(data, &result1)
	fmt.Printf("%v\n", result1)
	assert.Nil(t, err)

	data, err = disorder.Marshal(&result1)
	assert.Nil(t, err)
	result2 := test.Object{}
	err = disorder.Unmarshal(data, &result2)
	fmt.Printf("%v\n", result2)
	assert.Nil(t, err)
	assert.EqualValues(t, target, result2)

	json1, err := json.Marshal(result1)
	assert.Nil(t, err)
	json2, err := json.Marshal(result2)
	assert.Nil(t, err)

	assert.JSONEq(t, string(json1), string(json2))
}

func TestRpcMath(t *testing.T) {
	s := rpc.NewServer()
	sub.RegisterMathServiceServer(s, &testService{})
	err := s.Listen(":9999")
	assert.Nil(t, err)

	c := sub.NewMathServiceClient(rpc.NewClient("localhost:9999"))
	result, rpcErr := c.Increase(rpc.NewContext(), 17)

	assert.Nil(t, rpcErr)
	assert.Equal(t, int32(18), result)
}

func TestRpcPrimary(t *testing.T) {
	s := rpc.NewServer()
	test.RegisterPrimaryServiceServer(s, &testService{})
	_ = s.Listen(":8888")

	c := test.NewPrimaryServiceClient(rpc.NewClient("localhost:8888"))

	tt := time.Now()
	color := test.ColorRed
	result1, rpcErr := c.PrintObject(rpc.NewContext(), &test.Object{
		Time:        &tt,
		IntField:    123,
		StringField: "foo",
		EnumField:   &color,
		IntArray:    []int32{1, 2, 3},
		IntMap: map[string]int32{
			"1": 1,
			"2": 2,
			"3": 3,
		},
		ObjArray: []*sub.SubObject{{Value: 123}},
		ObjMap: map[string]*sub.SubObject{
			"foo": {Value: 123},
		},
	})
	assert.Nil(t, rpcErr)
	assert.Equal(t, int32(456), result1.IntField)
	assert.Equal(t, "bar", result1.StringField)
	newColor, err := result1.EnumField.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "green", newColor)
	assert.Equal(t, int32(4), result1.IntArray[0])
	assert.Equal(t, int32(4), result1.IntMap["4"])
	assert.Equal(t, int32(456), result1.ObjArray[0].Value)
	assert.Equal(t, int32(456), result1.ObjMap["bar"].Value)

	result2, rpcErr := c.PrintSubObject(rpc.NewContext(), &sub.SubObject{
		Value: 123,
	})
	assert.Nil(t, rpcErr)
	assert.Equal(t, int32(456), result2.Value)

	result3, rpcErr := c.PrintTime(rpc.NewContext(), &tt)
	assert.Nil(t, rpcErr)
	fmt.Printf("output time: %v", *result3)

	result4, rpcErr := c.PrintArray(rpc.NewContext(), []int32{1, 2, 3})
	assert.Nil(t, rpcErr)
	assert.Equal(t, 3, len(result4))
	assert.Equal(t, int32(4), result4[0])

	result5, rpcErr := c.PrintEnum(rpc.NewContext(), &color)
	assert.Nil(t, rpcErr)
	newColor, err = result5.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "green", newColor)

	result6, rpcErr := c.PrintMap(rpc.NewContext(), map[string]string{
		"foo": "bar",
	})
	assert.Nil(t, rpcErr)
	assert.Equal(t, "foo", result6["bar"])
}
