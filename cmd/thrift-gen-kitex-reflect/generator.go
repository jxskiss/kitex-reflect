package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cloudwego/thriftgo/parser"

	idl "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl"
)

type pluginGenerator struct {
	ast  *parser.Thrift
	desc *idl.ServiceDesc
}

func (p *pluginGenerator) getPkgName() string {
	return strings.ToLower(p.ast.Services[0].Name)
}

func (p *pluginGenerator) getOutputFile(outputPath string) string {
	var namespace string
	for _, ns := range p.ast.Namespaces {
		if ns.Language == "go" {
			namespace = strings.ReplaceAll(ns.Name, ".", "/")
		}
	}
	svcName := strings.ToLower(p.ast.Services[0].Name)
	filename := "plugin-reflect-desc.go"
	descFile := filepath.Join(outputPath, namespace, svcName, filename)
	return descFile
}

func (p *pluginGenerator) genServiceDesc() (string, error) {
	svc := p.ast.Services[0]
	p.desc = &idl.ServiceDesc{
		Name:         svc.Name,
		Functions:    nil,
		TypeDescList: nil,
	}
	for _, fn := range svc.Functions {
		fnDesc, err := p.getFunctionDesc(p.ast, fn)
		if err != nil {
			return "", fmt.Errorf("cannot get function descriptor: %w", err)
		}
		p.desc.Functions = append(p.desc.Functions, fnDesc)
	}
	jsonBuf, err := json.MarshalIndent(p.desc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("cannot marshal service descriptor: %w", err)
	}
	str := "`" + string(jsonBuf) + "`"
	return str, nil
}

func (p *pluginGenerator) getFunctionDesc(tree *parser.Thrift, fn *parser.Function) (*idl.FunctionDesc, error) {
	reqDesc, err := p.getReqTypeDesc(tree, fn.Arguments)
	if err != nil {
		return nil, err
	}
	respDesc, err := p.getRespTypeDesc(tree, fn.FunctionType)
	if err != nil {
		return nil, err
	}
	annotations := p.convAnnotations(fn.Annotations)
	hasRequestBase := false
	for _, reqField := range reqDesc.Struct.Fields[0].Type.Struct.Fields {
		if reqField.Type.GetIsRequestBase() {
			hasRequestBase = true
		}
	}
	desc := &idl.FunctionDesc{
		Name:           fn.Name,
		Oneway:         fn.Oneway,
		HasRequestBase: hasRequestBase,
		Request:        reqDesc,
		Response:       respDesc,
		Annotations:    annotations,
	}
	return desc, nil
}

func (p *pluginGenerator) getReqTypeDesc(tree *parser.Thrift, args []*parser.Field) (*idl.TypeDesc, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("unsupported arguments count: %d", len(args))
	}
	arg := args[0]
	if arg.Type.Category != parser.Category_Struct {
		return nil, fmt.Errorf("unsupported request type: %v", arg.Type.Category)
	}

	reqDesc := &idl.TypeDesc{
		Name: "",
		Type: idl.Type_STRUCT,
		Key:  nil,
		Elem: nil,
		Struct: &idl.StructDesc{
			Name:        "",
			Fields:      nil,
			Annotations: nil,
		},
		IsRequestBase: nil,
	}
	argDesc, err := p.getStructDesc(tree, arg.Type.Name, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot get request struct descriptor: %w", err)
	}
	fieldDesc := &idl.FieldDesc{
		Name:          arg.Name,
		Alias:         "",
		ID:            arg.ID,
		Required:      arg.Requiredness == parser.FieldType_Required,
		Optional:      arg.Requiredness == parser.FieldType_Optional,
		IsException:   false,
		DefaultValue:  nil,
		TypeDescIndex: nil,
		Type: &idl.TypeDesc{
			Name:          argDesc.Name,
			Type:          idl.Type_STRUCT,
			Key:           nil,
			Elem:          nil,
			Struct:        argDesc,
			IsRequestBase: nil,
		},
		Annotations: nil,
	}
	reqDesc.Struct.Fields = append(reqDesc.Struct.Fields, fieldDesc)
	return reqDesc, nil
}

func (p *pluginGenerator) getRespTypeDesc(tree *parser.Thrift, typ *parser.Type) (*idl.TypeDesc, error) {
	if typ.Category != parser.Category_Struct {
		return nil, fmt.Errorf("unsupported response type: %v", typ.Category)
	}
	structDesc, err := p.getStructDesc(tree, typ.Name, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot get response struct descriptor: %w", err)
	}
	desc := &idl.TypeDesc{
		Name:          typ.Name,
		Type:          idl.Type_STRUCT,
		Key:           nil,
		Elem:          nil,
		Struct:        structDesc,
		IsRequestBase: nil,
	}
	return desc, nil
}

func (p *pluginGenerator) getStructDesc(tree *parser.Thrift, structName string, recurDepth int) (*idl.StructDesc, error) {
	var structTyp *parser.StructLike
	var found bool
	typPkg, typName := splitType(structName)
	if typPkg != "" {
		ref, ok := tree.GetReference(typPkg)
		if !ok {
			return nil, fmt.Errorf("reference is missingf: %s", typPkg)
		}
		tree = ref
	}
	structTyp, found = findStructLike(tree, typName)
	if !found {
		return nil, fmt.Errorf("cannot find struct %v", structName)
	}
	annotations := p.convAnnotations(structTyp.Annotations)
	desc := &idl.StructDesc{
		Name:        typName,
		Fields:      nil,
		Annotations: annotations,
	}
	for _, field := range structTyp.Fields {
		fDesc, err := p.getFieldDesc(tree, field, recurDepth+1)
		if err != nil {
			return nil, fmt.Errorf("cannot get field descriptor: %w", err)
		}
		desc.Fields = append(desc.Fields, fDesc)
	}
	return desc, nil
}

func (p *pluginGenerator) getFieldDesc(tree *parser.Thrift, field *parser.Field, recurDepth int) (*idl.FieldDesc, error) {
	typDesc, err := p.getTypeDesc(tree, field.Type, recurDepth)
	if err != nil {
		return nil, fmt.Errorf("cannot get descriptor for type %s: %w", field.Type.Name, err)
	}
	annotations := p.convAnnotations(field.Annotations)
	fieldDesc := &idl.FieldDesc{
		Name:         field.Name,
		Alias:        "",
		ID:           field.ID,
		Required:     field.Requiredness == parser.FieldType_Required,
		Optional:     field.Requiredness == parser.FieldType_Optional,
		IsException:  false,
		DefaultValue: nil,
		Type:         typDesc,
		Annotations:  annotations,
	}
	if field.Default != nil {
		val, err := p.parseDefaultValue(tree, field.Name, field.Type, field.Default)
		if err != nil {
			return nil, fmt.Errorf("failed to parse default value: %w", err)
		}
		fieldDesc.DefaultValue = &val
	}
	return fieldDesc, nil
}

func (p *pluginGenerator) getTypeDesc(tree *parser.Thrift, typ *parser.Type, recurDepth int) (*idl.TypeDesc, error) {
	if typ == nil {
		return nil, nil
	}
	typEnum, err := p.convTypeEnum(typ)
	if err != nil {
		return nil, fmt.Errorf("cannot convert type enum: %w", err)
	}
	keyDesc, err := p.getTypeDesc(tree, typ.KeyType, recurDepth+1)
	if err != nil {
		return nil, fmt.Errorf("cannot get key type descriptor: %w", err)
	}
	elemDesc, err := p.getTypeDesc(tree, typ.ValueType, recurDepth+1)
	if err != nil {
		return nil, fmt.Errorf("cannot get elem type descriptor: %w", err)
	}

	var structDesc *idl.StructDesc
	var isRequestBase bool
	if typ.Category == parser.Category_Struct && typ.Name != "" {
		structDesc, err = p.getStructDesc(tree, typ.Name, recurDepth)
		if err != nil {
			return nil, fmt.Errorf("cannot get struct descriptor: %w", err)
		}
		isRequestBase = typ.Name == "base.Base" && recurDepth == 1
	}
	desc := &idl.TypeDesc{
		Name:          typ.Name,
		Type:          typEnum,
		Key:           keyDesc,
		Elem:          elemDesc,
		Struct:        structDesc,
		IsRequestBase: &isRequestBase,
	}
	return desc, nil
}

func (p *pluginGenerator) convAnnotations(annotations parser.Annotations) []*idl.Annotation {
	ret := make([]*idl.Annotation, 0, len(annotations))
	for _, ann := range annotations {
		idlAnn := &idl.Annotation{
			Key:    ann.Key,
			Values: ann.Values,
		}
		ret = append(ret, idlAnn)
	}
	return ret
}

var typeEnumTable = [...]idl.Type{
	parser.Category_Bool:      idl.Type_BOOL,
	parser.Category_Byte:      idl.Type_BYTE,
	parser.Category_I16:       idl.Type_I16,
	parser.Category_I32:       idl.Type_I32,
	parser.Category_I64:       idl.Type_I64,
	parser.Category_Double:    idl.Type_DOUBLE,
	parser.Category_String:    idl.Type_STRING,
	parser.Category_Binary:    idl.Type_STRING,
	parser.Category_Map:       idl.Type_MAP,
	parser.Category_List:      idl.Type_LIST,
	parser.Category_Set:       idl.Type_SET,
	parser.Category_Enum:      idl.Type_I32,
	parser.Category_Struct:    idl.Type_STRUCT,
	parser.Category_Union:     idl.Type_STRUCT,
	parser.Category_Exception: idl.Type_STRUCT,
}

func (p *pluginGenerator) convTypeEnum(typ *parser.Type) (idl.Type, error) {
	var idlTyp idl.Type
	if int(typ.Category) < len(typeEnumTable) {
		idlTyp = typeEnumTable[typ.Category]
	}
	if idlTyp > 0 {
		return idlTyp, nil
	}

	if typ.Category == parser.Category_Typedef {
		underlyingTyp, found := p.ast.GetTypedef(typ.Name)
		if !found {
			return 0, fmt.Errorf("cannot find typedef %s", typ.Name)
		}
		return p.convTypeEnum(underlyingTyp.Type)
	}

	return 0, fmt.Errorf("unsupported type category: %v", typ.Category)
}

func (p *pluginGenerator) parseDefaultValue(tree *parser.Thrift, name string, t *parser.Type, v *parser.ConstValue) (string, error) {
	val, err := parse(tree, name, t, v)
	if err != nil {
		return "", err
	}
	jBuf, err := json.Marshal(val)
	if err != nil {
		return "", err
	}
	return string(jBuf), nil
}
