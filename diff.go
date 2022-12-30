package kitexreflect

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/cloudwego/kitex/pkg/generic/descriptor"
)

func NewDiffComparator(kitexDesc, pluginDesc *descriptor.ServiceDescriptor) *DiffComparator {
	cmp := &DiffComparator{
		kDesc: kitexDesc,
		pDesc: pluginDesc,
	}
	return cmp
}

type DiffComparator struct {
	kDesc    *descriptor.ServiceDescriptor
	pDesc    *descriptor.ServiceDescriptor
	messages []string
}

func (cmp *DiffComparator) Messages() []string {
	return cmp.messages
}

func (cmp *DiffComparator) ErrorMessages() []string {
	var out []string
	for _, msg := range cmp.messages {
		if strings.HasPrefix(msg, "ERROR:") {
			out = append(out, msg)
		}
	}
	return out
}

func (cmp *DiffComparator) InfoMessages() []string {
	var out []string
	for _, msg := range cmp.messages {
		if strings.HasPrefix(msg, "INFO:") {
			out = append(out, msg)
		}
	}
	return out
}

func (cmp *DiffComparator) Check() error {
	cmp.diffService(cmp.kDesc, cmp.pDesc)
	return nil
}

func (cmp *DiffComparator) diffService(k, p *descriptor.ServiceDescriptor) {
	cmp.checkEqual("svc", "serviceName", k.Name, p.Name)
	cmp.checkEqual("svc", "functions count", len(k.Functions), len(p.Functions))

	kitexFuncNames := getSortedKeys(k.Functions)
	pluginFuncNames := getSortedKeys(p.Functions)
	cmp.checkEqual("svc", "function names", kitexFuncNames, pluginFuncNames)

	for i, name := range kitexFuncNames {
		cmp.log(fmt.Sprintf("INFO: [svc] checking service method: %d, %s", i, name))
		kitexFunc, err := k.LookupFunctionByMethod(name)
		if err != nil {
			cmp.log(fmt.Sprintf("ERROR: [svc] cannot find function %s from kitex's service descriptor", name))
			continue
		}
		pluginFunc, err := p.LookupFunctionByMethod(name)
		if err != nil {
			cmp.log(fmt.Sprintf("ERROR: [svc] cannot find function %s from plugin's service descriptor", name))
			continue
		}
		cmp.diffFunction(kitexFunc, pluginFunc)
	}

	cmp.diffRouter(k, p)
}

func (cmp *DiffComparator) diffRouter(k, p *descriptor.ServiceDescriptor) {
	if k.Router == nil && p.Router == nil {
		cmp.log("INFO: [router] both router is nil")
		return
	}
	kr := reflect.ValueOf(k.Router)
	pr := reflect.ValueOf(p.Router)
	kTrees := kr.Elem().FieldByName("trees")
	pTrees := pr.Elem().FieldByName("trees")
	cmp.checkEqual("router", "router length", kTrees.Len(), pTrees.Len())

	for _, key := range kTrees.MapKeys() {
		kNode := kTrees.MapIndex(key)
		pNode := pTrees.MapIndex(key)
		cmp.diffRouterNode("router.node", kNode, pNode)
	}
}

func (cmp *DiffComparator) diffRouterNode(prefix string, k, p reflect.Value) {
	cmp.checkEqual(prefix, "node is valid", k.IsValid(), p.IsValid())
	if !k.IsValid() {
		return
	}

	cmp.checkEqual(prefix, "node path",
		k.Elem().FieldByName("path").String(),
		p.Elem().FieldByName("path").String(),
	)

	kChildren := k.Elem().FieldByName("children")
	pChildren := p.Elem().FieldByName("children")
	cmp.checkEqual(prefix, "node children length", kChildren.Len(), pChildren.Len())

	for i := 0; i < kChildren.Len(); i++ {
		nextPrefix := fmt.Sprintf("%s.children[%d]", prefix, i)
		kNode := kChildren.Index(i)
		pNode := pChildren.Index(i)
		cmp.diffRouterNode(nextPrefix, kNode, pNode)
	}
}

func (cmp *DiffComparator) diffFunction(k, p *descriptor.FunctionDescriptor) {
	cmp.checkEqual("svc", "function name", k.Name, p.Name)
	cmp.checkEqual(k.Name, "function oneway", k.Oneway, p.Oneway)
	cmp.checkEqual(k.Name, "function hasRequestBase", k.HasRequestBase, p.HasRequestBase)

	cmp.diffFuncRequest(k.Name, k.Request, p.Request)
	cmp.diffFuncResponse(k.Name, k.Response, p.Response)
}

func (cmp *DiffComparator) diffFuncRequest(prefix string, k, p *descriptor.TypeDescriptor) {
	cmp.checkEqual(prefix, "request name", k.Name, p.Name)
	cmp.checkEqual(prefix, "request type", k.Type, p.Type)
	cmp.checkEqual(prefix, "request isRequestBase", k.IsRequestBase, p.IsRequestBase)

	cmp.diffTypeDescriptor(prefix+".req.key", k.Key, p.Key)
	cmp.diffTypeDescriptor(prefix+".req.elem", k.Elem, p.Elem)
	cmp.diffStructDescriptor(prefix+".req.struct", k.Struct, p.Struct)
}

func (cmp *DiffComparator) diffFuncResponse(prefix string, k, p *descriptor.TypeDescriptor) {
	cmp.checkEqual(prefix, "response name", k.Name, p.Name)
	cmp.checkEqual(prefix, "response type", k.Type, p.Type)
	cmp.checkEqual(prefix, "response isRequestBase", k.IsRequestBase, p.IsRequestBase)

	cmp.diffTypeDescriptor(prefix+".resp.key", k.Key, p.Key)
	cmp.diffTypeDescriptor(prefix+".resp.elem", k.Elem, p.Elem)
	cmp.diffStructDescriptor(prefix+".resp.struct", k.Struct, p.Struct)
}

func (cmp *DiffComparator) diffTypeDescriptor(prefix string, k, p *descriptor.TypeDescriptor) {
	if k == nil && p == nil {
		cmp.log(fmt.Sprintf("INFO: [%s] both type desc is nil", prefix))
		return
	}

	cmp.checkEqual(prefix, "type name", k.Name, p.Name)
	cmp.checkEqual(prefix, "type type", k.Type, p.Type)
	cmp.checkEqual(prefix, "type isRequestBase", k.IsRequestBase, p.IsRequestBase)

	cmp.diffTypeDescriptor(prefix+".key", k.Key, p.Key)
	cmp.diffTypeDescriptor(prefix+".elem", k.Elem, p.Elem)
	cmp.diffStructDescriptor(prefix+".struct", k.Struct, p.Struct)
}

func (cmp *DiffComparator) diffStructDescriptor(prefix string, k, p *descriptor.StructDescriptor) {
	if k == nil && p == nil {
		cmp.log(fmt.Sprintf("INFO: [%s] both struct desc is nil", prefix))
		return
	}

	cmp.checkEqual(prefix, "struct name", k.Name, p.Name)
	cmp.checkEqual(prefix, "struct field IDs", getSortedKeys(k.FieldsByID), getSortedKeys(p.FieldsByID))
	cmp.checkEqual(prefix, "struct field names", getSortedKeys(k.FieldsByName), getSortedKeys(p.FieldsByName))
	cmp.checkEqual(prefix, "struct required fields", getSortedKeys(k.RequiredFields), getSortedKeys(p.RequiredFields))
	cmp.checkEqual(prefix, "struct default fields", getSortedKeys(k.DefaultFields), getSortedKeys(p.DefaultFields))

	fieldIDs := getSortedKeys(k.FieldsByID)
	for _, id := range fieldIDs {
		kField := k.FieldsByID[id]
		pField := p.FieldsByID[id]
		cmp.diffStructField(prefix, kField, pField)
	}
}

func (cmp *DiffComparator) diffStructField(prefix string, k, p *descriptor.FieldDescriptor) {
	cmp.checkEqual(prefix, "field id", k.ID, p.ID)
	cmp.checkEqual(prefix, "field name", k.Name, p.Name)
	cmp.checkEqual(prefix, "field alias", k.Alias, p.Alias)
	cmp.checkEqual(prefix, "field isRequired", k.Required, p.Required)
	cmp.checkEqual(prefix, "field isOptional", k.Optional, p.Optional)
	cmp.checkEqual(prefix, "field isException", k.IsException, p.IsException)
	cmp.checkEqual(prefix, "field default value", k.DefaultValue, p.DefaultValue)
	cmp.checkEqual(prefix, "field has http mapping", k.HTTPMapping != nil, p.HTTPMapping != nil)
	cmp.checkEqual(prefix, "field has value mapping", k.ValueMapping != nil, p.ValueMapping != nil)

	cmp.diffTypeDescriptor(prefix+"."+k.Name, k.Type, p.Type)
}

func (cmp *DiffComparator) checkEqual(prefix, msg string, x, y interface{}) {
	var cmpMsg string
	equal := reflect.DeepEqual(x, y)
	if !equal {
		cmpMsg = fmt.Sprintf("ERROR: [%s] %s not equal, %v != %v", prefix, msg, x, y)
	} else {
		cmpMsg = fmt.Sprintf("INFO: [%s] %s equal, %v == %v", prefix, msg, x, y)
	}
	cmp.log(cmpMsg)
}

func (cmp *DiffComparator) log(msg string) {
	cmp.messages = append(cmp.messages, msg)
}

type sortable interface {
	int | int32 | int64 | string
}

func getSortedKeys[M ~map[K]V, K sortable, V any](m M) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}
