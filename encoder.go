package disorder

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"time"
)

type Encoder struct {
	writer io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		writer: w,
	}
}

func (e *Encoder) Encode(value interface{}) error {
	return e.write(reflect.ValueOf(value))
}

func (e *Encoder) write(value reflect.Value) error {
	if !value.IsValid() || (value.Kind() == reflect.Ptr && value.IsNil()) {
		return fmt.Errorf("invalid value or type")
	}

	switch i := value.Interface().(type) {
	case time.Time:
		return e.writeTime(&i)

	case *time.Time:
		return e.writeTime(i)

	case Enum:
		_, err := e.writer.Write([]byte{byte(tagEnum)})
		if err != nil {
			return err
		}
		name, err := i.ToString()
		if err != nil {
			return err
		}
		return e.writeName(name)

	case Marshaler:
		return i.MarshalDO(e.writer)
	}

	var bytes []byte
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		return e.write(value.Elem())

	case reflect.Bool:
		bytes = make([]byte, 2)
		bytes[0] = byte(tagBool)
		if value.Bool() {
			bytes[1] = 1
		} else {
			bytes[1] = 0
		}

	case reflect.Uint8:
		bytes = make([]byte, 2)
		bytes[0] = byte(tagU8)
		bytes[1] = byte(uint8(value.Uint()))

	case reflect.Uint16:
		bytes = make([]byte, 3)
		bytes[0] = byte(tagU16)
		binary.LittleEndian.PutUint16(bytes[1:], uint16(value.Uint()))

	case reflect.Uint32, reflect.Uint:
		bytes = make([]byte, 5)
		bytes[0] = byte(tagU32)
		binary.LittleEndian.PutUint32(bytes[1:], uint32(value.Uint()))

	case reflect.Uint64:
		bytes = make([]byte, 9)
		bytes[0] = byte(tagU64)
		binary.LittleEndian.PutUint64(bytes[1:], uint64(value.Uint()))

	case reflect.Int8:
		bytes = make([]byte, 2)
		bytes[0] = byte(tagI8)
		bytes[1] = byte(int8(value.Int()))

	case reflect.Int16:
		bytes = make([]byte, 3)
		bytes[0] = byte(tagI16)
		binary.LittleEndian.PutUint16(bytes[1:], uint16(value.Int()))

	case reflect.Int32, reflect.Int:
		bytes = make([]byte, 5)
		bytes[0] = byte(tagI32)
		binary.LittleEndian.PutUint32(bytes[1:], uint32(value.Int()))

	case reflect.Int64:
		bytes = make([]byte, 9)
		bytes[0] = byte(tagI64)
		binary.LittleEndian.PutUint64(bytes[1:], uint64(value.Int()))

	case reflect.Float32:
		bytes = make([]byte, 5)
		bytes[0] = byte(tagF32)
		binary.LittleEndian.PutUint32(bytes[1:], math.Float32bits(float32(value.Float())))

	case reflect.Float64:
		bytes = make([]byte, 9)
		bytes[0] = byte(tagF64)
		binary.LittleEndian.PutUint64(bytes[1:], math.Float64bits(value.Float()))

	case reflect.String:
		bytes = make([]byte, 5)
		bytes[0] = byte(tagString)
		str := value.String()
		binary.LittleEndian.PutUint32(bytes[1:], uint32(len(str)))
		bytes = append(bytes, []byte(str)...)

	case reflect.Slice, reflect.Array:
		return e.writeArray(value)

	case reflect.Map:
		return e.writeMap(value)

	case reflect.Struct:
		return e.writeObject(value)

	default:
		return fmt.Errorf("invalid type: %s", value.Type().String())
	}

	_, err := e.writer.Write(bytes)
	return err
}

func (e *Encoder) writeArray(value reflect.Value) error {
	count := value.Len()
	if count == 0 {
		return nil
	}

	_, err := e.writer.Write([]byte{byte(tagStartArray)})
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		err = e.write(value.Index(i))
		if err != nil {
			return err
		}
	}

	_, err = e.writer.Write([]byte{byte(tagEndArray)})
	return err
}

func (e *Encoder) writeMap(value reflect.Value) error {
	count := value.Len()
	if count == 0 {
		return nil
	}

	_, err := e.writer.Write([]byte{byte(tagStartObject)})
	if err != nil {
		return err
	}

	keys := value.MapKeys()
	for _, key := range keys {
		if key.Kind() != reflect.String {
			return fmt.Errorf("map key type must be string")
		}
		err := e.writeName(key.String())
		if err != nil {
			return err
		}
		err = e.write(value.MapIndex(key))
		if err != nil {
			return err
		}
	}

	_, err = e.writer.Write([]byte{byte(tagEndObject)})
	return err
}

func (e *Encoder) writeObject(value reflect.Value) error {
	info, err := getStructInfo(value.Type())
	if err != nil {
		return err
	}

	_, err = e.writer.Write([]byte{byte(tagStartObject)})
	if err != nil {
		return err
	}

	for _, field := range info.fieldsList {
		fieldValue := value.Field(field.index)
		if field.omitEmpty && isZero(fieldValue) {
			continue
		}
		err = e.writeName(field.key)
		if err != nil {
			return err
		}
		err = e.write(fieldValue)
		if err != nil {
			return err
		}
	}

	_, err = e.writer.Write([]byte{byte(tagEndObject)})
	return err
}

func (e *Encoder) writeTime(t *time.Time) error {
	bytes := make([]byte, 9)
	bytes[0] = byte(tagTimestamp)
	binary.LittleEndian.PutUint64(bytes[1:], uint64(t.Unix()))
	_, err := e.writer.Write(bytes)
	return err
}

func (e *Encoder) writeName(value string) error {
	if len(value) > 255 {
		return fmt.Errorf("string length overflow. should less than 255")
	}
	_, err := e.writer.Write([]byte{byte(len(value))})
	if err != nil {
		return err
	}
	_, err = e.writer.Write([]byte(value))
	return err
}
