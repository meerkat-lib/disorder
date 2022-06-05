package loader

import (
	"fmt"

	"github.com/meerkat-lib/disorder/internal/schema"
)

var (
	undefined = &schema.TypeInfo{}
)

type parser struct {
	validator *validator
}

func newParser() *parser {
	return &parser{
		validator: newValidator(),
	}
}

func (p *parser) qualifiedName(pkg, name string) string {
	return fmt.Sprintf("%s.%s", pkg, name)
}

func (p *parser) parse(proto *proto) (*schema.File, error) {
	if proto.Package == "" {
		return nil, fmt.Errorf("package name is required")
	}
	if !p.validator.validatePackageName(proto.Package) {
		return nil, fmt.Errorf("invalid package name: %s", proto.Package)
	}
	file := &schema.File{
		FilePath: proto.FilePath,
		Package:  proto.Package,
	}
	for name, values := range proto.Enums {
		if !p.validator.validateEnumName(name) {
			return nil, fmt.Errorf("invalid enum name: %s", name)
		}
		valuesSet := map[string]bool{}
		enum := &schema.Enum{
			Name: p.qualifiedName(proto.Package, name),
		}
		for _, value := range values {
			if _, exists := valuesSet[value]; exists {
				return nil, fmt.Errorf("duplicated enum value [%s]", value)
			}
			if !p.validator.validateEnumValue(value) {
				return nil, fmt.Errorf("invalid enum value: %s", value)
			}
			valuesSet[value] = true
			enum.Values = append(enum.Values, value)
		}
		file.Enums = append(file.Enums, enum)
	}
	for name, fields := range proto.Messages {
		if !p.validator.validateMessageName(name) {
			return nil, fmt.Errorf("invalid message name: %s", name)
		}
		fieldsSet := map[string]bool{}
		message := &schema.Message{
			Name: p.qualifiedName(proto.Package, name),
		}
		for fieldName, fieldType := range fields {
			if _, exists := fieldsSet[fieldName]; exists {
				return nil, fmt.Errorf("duplicated field [%s]", fieldName)
			}
			if !p.validator.validateFieldName(fieldName) {
				return nil, fmt.Errorf("invalid field name: %s", fieldName)
			}
			fieldsSet[fieldName] = true
			field, err := p.parseField(proto.Package, fieldName, fieldType)
			if err != nil {
				return nil, err
			}
			message.Fields = append(message.Fields, field)
		}
		file.Messages = append(file.Messages, message)
	}
	for name, rpcs := range proto.Services {
		if !p.validator.validateServiceName(name) {
			return nil, fmt.Errorf("invalid service name: %s", name)
		}
		rpcsSet := map[string]bool{}
		service := &schema.Service{
			Name: name,
		}
		for rpcName, rpcDefine := range rpcs {
			if _, exists := rpcsSet[rpcName]; exists {
				return nil, fmt.Errorf("duplicated rpc [%s]", rpcName)
			}
			if !p.validator.validateRpcName(rpcName) {
				return nil, fmt.Errorf("invalid rpc name: %s", rpcName)
			}
			rpcsSet[rpcName] = true
			rpc, err := p.parseRpc(proto.Package, rpcName, rpcDefine)
			if err != nil {
				return nil, err
			}
			service.Rpc = append(service.Rpc, rpc)
		}
		file.Services = append(file.Services, service)
	}
	return file, nil
}

func (p *parser) parseField(pkg, name, typ string) (*schema.Field, error) {
	info, err := p.parseType(pkg, typ)
	if err != nil {
		return nil, fmt.Errorf("field [%s] error: %s", name, err.Error())
	}
	return &schema.Field{
		Name: name,
		Type: info,
	}, nil
}

func (p *parser) parseRpc(pkg, name string, rpc *rpc) (*schema.Rpc, error) {
	var err error
	r := &schema.Rpc{
		Name: name,
	}
	if rpc.Input == "" {
		r.Input = undefined
	} else {
		r.Input, err = p.parseType(pkg, rpc.Input)
		if err != nil {
			return nil, fmt.Errorf("rpc [%s] input type error: %s", name, err.Error())
		}
	}
	if rpc.Output == "" {
		r.Output = undefined
	} else {
		r.Output, err = p.parseType(pkg, rpc.Output)
		if err != nil {
			return nil, fmt.Errorf("rpc [%s] output type error: %s", name, err.Error())
		}
	}
	return r, nil
}

func (p *parser) parseType(pkg, typ string) (t *schema.TypeInfo, err error) {
	if typ == "" {
		err = fmt.Errorf("empty type")
		return
	}
	t = &schema.TypeInfo{}
	if p.validator.isSimpleType(typ) {
		if p.validator.isPrimary(typ) {
			t.Type = p.validator.primaryType(typ)
			return
		} else {
			t.TypeRef = p.qualifiedName(pkg, typ)
			return
		}
	} else if p.validator.isQualifiedType(typ) {
		t.TypeRef = p.qualifiedName(pkg, typ)
		return
	} else if p.validator.isSimpleArrayType(typ) {
		t.Type = schema.TypeArray
		subType := typ[6 : len(typ)-1]
		if p.validator.isPrimary(subType) {
			t.SubType = p.validator.primaryType(subType)
			return
		} else {
			t.SubTypeRef = p.qualifiedName(pkg, subType)
			return
		}
	} else if p.validator.isQualifiedArrayType(typ) {
		t.Type = schema.TypeArray
		t.SubTypeRef = typ[6 : len(typ)-1]
		return
	} else if p.validator.isSimpleMapType(typ) {
		t.Type = schema.TypeMap
		subType := typ[4 : len(typ)-1]
		if p.validator.isPrimary(subType) {
			t.SubType = p.validator.primaryType(subType)
			return
		} else {
			t.SubTypeRef = p.qualifiedName(pkg, subType)
			return
		}
	} else if p.validator.isQualifiedMapType(typ) {
		t.Type = schema.TypeMap
		t.SubTypeRef = typ[4 : len(typ)-1]
		return
	}
	return nil, fmt.Errorf("invalid type %s", typ)
}
