package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cloudwego/thriftgo/parser"

	idl "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl"
)

func newPluginGenerator(ast *parser.Thrift) *pluginGenerator {
	gen := &pluginGenerator{
		ast:         ast,
		namespace:   make(map[string]int),
		structIndex: make(map[string]int),
	}
	return gen
}

type pluginGenerator struct {
	ast  *parser.Thrift
	desc *idl.ServiceDesc

	namespace map[string]int

	structList  []*idl.StructDesc
	structIndex map[string]int
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

func (p *pluginGenerator) buildServiceDesc() (string, error) {
	svc := p.ast.Services[0]
	p.desc = &idl.ServiceDesc{
		Name:      svc.Name,
		Functions: nil,
	}
	for _, fn := range svc.Functions {
		fnDesc, err := p.buildFunctionDesc(p.ast, fn)
		if err != nil {
			return "", fmt.Errorf("cannot get function descriptor: %w", err)
		}
		p.desc.Functions = append(p.desc.Functions, fnDesc)
	}
	p.desc.StructList = p.structList
	jsonBuf, err := json.MarshalIndent(p.desc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("cannot marshal service descriptor: %w", err)
	}
	str := "`" + string(jsonBuf) + "`"
	return str, nil
}

func (p *pluginGenerator) buildFunctionDesc(tree *parser.Thrift, fn *parser.Function) (*idl.FunctionDesc, error) {
	reqDesc, err := p.buildReqTypeDesc(tree, fn.Arguments)
	if err != nil {
		return nil, err
	}
	respDesc, err := p.buildRespTypeDesc(tree, fn.FunctionType)
	if err != nil {
		return nil, err
	}
	annotations := p.convAnnotations(fn.Annotations)
	hasRequestBase := false
	reqStructDesc := p.getStructDesc(reqDesc.Struct.Fields[0].Type)
	for _, reqField := range reqStructDesc.Fields {
		if reqField.Type.GetIsRequestBase() {
			hasRequestBase = true
		}
	}
	desc := &idl.FunctionDesc{
		Name:           &fn.Name,
		Oneway:         nil,
		HasRequestBase: nil,
		Request:        reqDesc,
		Response:       respDesc,
		Annotations:    annotations,
	}
	setNonZeroValue(&desc.Oneway, fn.Oneway)
	setNonZeroValue(&desc.HasRequestBase, hasRequestBase)
	return desc, nil
}

func (p *pluginGenerator) buildReqTypeDesc(tree *parser.Thrift, args []*parser.Field) (*idl.TypeDesc, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("unsupported arguments count: %d", len(args))
	}
	arg := args[0]
	if arg.Type.Category != parser.Category_Struct {
		return nil, fmt.Errorf("unsupported request type: %v", arg.Type.Category)
	}

	reqDesc := &idl.TypeDesc{
		Name: nil,
		Type: idlTypeStruct,
		Struct: &idl.StructDesc{
			Name:        nil,
			Fields:      nil,
			Annotations: nil,
		},
	}
	structIdx, err := p.buildStructDesc(tree, arg.Type.Name, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot get request struct descriptor: %w", err)
	}
	fieldDesc := &idl.FieldDesc{
		Name:         &arg.Name,
		Alias:        nil,
		ID:           &arg.ID,
		Required:     nil,
		Optional:     nil,
		IsException:  nil,
		DefaultValue: nil,
		Type: &idl.TypeDesc{
			Name:      &arg.Type.Name,
			Type:      idlTypeStruct,
			StructIdx: int32p(structIdx),
		},
		Annotations: nil,
	}
	setNonZeroValue(&fieldDesc.Required, arg.Requiredness == parser.FieldType_Required)
	setNonZeroValue(&fieldDesc.Optional, arg.Requiredness == parser.FieldType_Optional)
	reqDesc.Struct.Fields = append(reqDesc.Struct.Fields, fieldDesc)
	return reqDesc, nil
}

func (p *pluginGenerator) buildRespTypeDesc(tree *parser.Thrift, typ *parser.Type) (*idl.TypeDesc, error) {
	if typ.Category != parser.Category_Struct {
		return nil, fmt.Errorf("unsupported response type: %v", typ.Category)
	}

	respDesc := &idl.TypeDesc{
		Name: nil,
		Type: idlTypeStruct,
		Struct: &idl.StructDesc{
			Name:        nil,
			Fields:      nil,
			Annotations: nil,
		},
	}
	structIdx, err := p.buildStructDesc(tree, typ.Name, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot get response struct descriptor: %w", err)
	}
	fieldDesc := &idl.FieldDesc{
		Type: &idl.TypeDesc{
			Name:      &typ.Name,
			Type:      idlTypeStruct,
			StructIdx: int32p(structIdx),
		},
		Annotations: nil,
	}
	respDesc.Struct.Fields = append(respDesc.Struct.Fields, fieldDesc)
	return respDesc, nil
}

func (p *pluginGenerator) buildStructDesc(tree *parser.Thrift, structName string, recurDepth int) (int, error) {
	var structTyp *parser.StructLike
	var found bool
	typPkg, typName := splitType(structName)
	if typPkg != "" {
		ref, ok := tree.GetReference(typPkg)
		if !ok {
			return -1, fmt.Errorf("reference is missingf: %s", typPkg)
		}
		tree = ref
	}

	typeNs := fmt.Sprintf("%s|%s", tree.Filename, typPkg)
	idxKey := fmt.Sprintf("%d|%s", p.getNamespaceID(typeNs), typName)
	structIdx, ok := p.structIndex[idxKey]
	if ok {
		return structIdx, nil
	}

	structTyp, found = findStructLike(tree, typName)
	if !found {
		return -1, fmt.Errorf("cannot find struct %v", structName)
	}
	annotations := p.convAnnotations(structTyp.Annotations)
	desc := &idl.StructDesc{
		Name:        &typName,
		Fields:      nil,
		Annotations: annotations,
		UniqueKey:   &idxKey,
	}

	// We need to add the struct to structList before fully resolve it
	// to avoid stack overflow when processing recursive types.
	structIdx = len(p.structList)
	p.structList = append(p.structList, desc)
	p.structIndex[idxKey] = structIdx

	for _, field := range structTyp.Fields {
		fDesc, err := p.buildFieldDesc(tree, field, recurDepth+1)
		if err != nil {
			return -1, fmt.Errorf("cannot get field descriptor: %w", err)
		}
		desc.Fields = append(desc.Fields, fDesc)
	}
	return structIdx, nil
}

func (p *pluginGenerator) getNamespaceID(ns string) int {
	idx, ok := p.namespace[ns]
	if !ok {
		idx = len(p.namespace)
		p.namespace[ns] = idx
	}
	return idx
}

func (p *pluginGenerator) getStructDesc(typeDesc *idl.TypeDesc) *idl.StructDesc {
	if idx := typeDesc.StructIdx; idx != nil {
		return p.structList[*idx]
	}
	return typeDesc.Struct
}

func (p *pluginGenerator) buildFieldDesc(tree *parser.Thrift, field *parser.Field, recurDepth int) (*idl.FieldDesc, error) {
	typDesc, err := p.buildTypeDesc(tree, field.Type, recurDepth)
	if err != nil {
		return nil, fmt.Errorf("cannot get descriptor for type %s: %w", field.Type.Name, err)
	}
	annotations := p.convAnnotations(field.Annotations)
	fieldDesc := &idl.FieldDesc{
		Name:         &field.Name,
		Alias:        nil,
		ID:           &field.ID,
		Required:     nil,
		Optional:     nil,
		IsException:  nil,
		DefaultValue: nil,
		Type:         typDesc,
		Annotations:  annotations,
	}
	setNonZeroValue(&fieldDesc.Required, field.Requiredness == parser.FieldType_Required)
	setNonZeroValue(&fieldDesc.Optional, field.Requiredness == parser.FieldType_Optional)
	if field.Default != nil {
		val, err := p.parseDefaultValue(tree, field.Name, field.Type, field.Default)
		if err != nil {
			return nil, fmt.Errorf("failed to parse default value: %w", err)
		}
		fieldDesc.DefaultValue = &val
	}
	return fieldDesc, nil
}

func (p *pluginGenerator) buildTypeDesc(tree *parser.Thrift, typ *parser.Type, recurDepth int) (*idl.TypeDesc, error) {
	if typ == nil {
		return nil, nil
	}

	if typDef, ok := tree.GetTypedef(typ.Name); ok {
		return p.buildTypeDesc(tree, typDef.Type, recurDepth+1)
	}

	typEnum, err := p.convTypeEnum(typ)
	if err != nil {
		return nil, fmt.Errorf("cannot convert type enum: %w", err)
	}

	desc := &idl.TypeDesc{
		Name:          &typ.Name,
		Type:          &typEnum,
		Key:           nil,
		Elem:          nil,
		Struct:        nil,
		IsRequestBase: nil,
	}
	if typ.Category == parser.Category_Enum {
		desc.Name = &enumTypeName
	}
	if typ.KeyType != nil {
		kbt := p.getBuiltinType(typ.KeyType)
		if kbt != idl.BuiltinType_NOT_BUITIN {
			desc.Kbt = &kbt
		} else {
			keyDesc, err := p.buildTypeDesc(tree, typ.KeyType, recurDepth+1)
			if err != nil {
				return nil, fmt.Errorf("cannot get key type descriptor: %w", err)
			}
			desc.Key = keyDesc
		}
	}
	if typ.ValueType != nil {
		ebt := p.getBuiltinType(typ.ValueType)
		if ebt != idl.BuiltinType_NOT_BUITIN {
			desc.Ebt = &ebt
		} else {
			elemDesc, err := p.buildTypeDesc(tree, typ.ValueType, recurDepth+1)
			if err != nil {
				return nil, fmt.Errorf("cannot get elem type descriptor: %w", err)
			}
			desc.Elem = elemDesc
		}
	}

	var structIdx = -1
	if typ.Category == parser.Category_Struct && typ.Name != "" {
		structIdx, err = p.buildStructDesc(tree, typ.Name, recurDepth)
		if err != nil {
			return nil, fmt.Errorf("cannot get struct descriptor: %w", err)
		}
		isRequestBase := typ.Name == "base.Base" && recurDepth == 1
		if isRequestBase {
			desc.IsRequestBase = &trueVal
		}
		if structIdx >= 0 {
			desc.StructIdx = int32p(structIdx)
		}
	}

	return desc, nil
}

var builtinTypeTable = map[string]idl.BuiltinType{
	"void":   idl.BuiltinType_VOID,
	"bool":   idl.BuiltinType_BOOL,
	"byte":   idl.BuiltinType_BYTE,
	"i8":     idl.BuiltinType_I8,
	"i16":    idl.BuiltinType_I16,
	"i32":    idl.BuiltinType_I32,
	"i64":    idl.BuiltinType_I64,
	"double": idl.BuiltinType_DOUBLE,
	"string": idl.BuiltinType_STRING,
	"binary": idl.BuiltinType_BINARY,
}

func (p *pluginGenerator) getBuiltinType(typ *parser.Type) idl.BuiltinType {
	return builtinTypeTable[typ.Name]
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
