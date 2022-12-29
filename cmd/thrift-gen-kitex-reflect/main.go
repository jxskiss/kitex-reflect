package main

import (
	"bytes"
	"encoding/json"
	"go/format"
	"io"
	"os"
	"text/template"
	"time"

	"github.com/cloudwego/thriftgo/plugin"
)

var _ json.Marshaler

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		println("Failed to read plugin input:", err.Error())
		os.Exit(1)
	}
	request, err := plugin.UnmarshalRequest(input)
	if err != nil {
		println("Failed to unmarshal plugin request:", err.Error())
		os.Exit(1)
	}
	if request.Language != "go" {
		println("Unsupported language:", request.Language)
		os.Exit(1)
	}

	debugln(func() {
		requestJSON, err := json.MarshalIndent(request, "", "  ")
		if err != nil {
			panic(err)
		}
		println(string(requestJSON))
	})

	gen := newPluginGenerator(request.AST)
	serviceDesc, err := gen.buildServiceDesc()
	if err != nil {
		println("Failed to generate service descriptor:", err.Error())
		os.Exit(1)
	}

	buf := &bytes.Buffer{}
	tpl, err := template.New("").Parse(implTpl)
	if err != nil {
		println("Failed to parse plugin template:", err.Error())
		os.Exit(1)
	}
	const genTimeLayout = "20060102150405"
	err = tpl.Execute(buf, map[string]interface{}{
		"PkgName":        gen.getPkgName(),
		"GenTime":        time.Now().Format(genTimeLayout),
		"GenServiceDesc": serviceDesc,
	})
	if err != nil {
		println("Failed to execute plugin template:", err.Error())
		os.Exit(1)
	}
	code, err := format.Source(buf.Bytes())
	if err != nil {
		println("Failed to format generated code:", err.Error())
		os.Exit(1)
	}

	outputFile := gen.getOutputFile(request.OutputPath)
	descCode := &plugin.Generated{
		Content:        string(code),
		Name:           &outputFile,
		InsertionPoint: nil,
	}
	response := &plugin.Response{
		Error:    nil,
		Contents: []*plugin.Generated{descCode},
		Warnings: nil,
	}
	os.Exit(exit(response))
}

func exit(response *plugin.Response) int {
	data, err := plugin.MarshalResponse(response)
	if err != nil {
		println("Failed to marshal response:", err.Error())
		return 1
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		println("Error at writing response out:", err.Error())
		return 1
	}
	return 0
}
