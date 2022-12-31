package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/cloudwego/thriftgo/generator/golang/extension/meta"
	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/plugin"
	"github.com/cloudwego/thriftgo/semantic"
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
	err = semantic.ResolveSymbols(request.AST)
	if err != nil {
		println("Failed to resolve symbols:", err.Error())
		os.Exit(1)
	}

	debugln(func() {
		requestJSON, err := json.MarshalIndent(request, "", "  ")
		if err != nil {
			panic(err)
		}
		println(string(requestJSON))
	})

	lastService := request.AST.Services[len(request.AST.Services)-1]
	pkgName := strings.ToLower(lastService.Name)

	idlBytes, err := encodeIDLToBytes(request.AST)
	if err != nil {
		println("Failed to encode IDL:", err.Error())
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
		"PkgName":  pkgName,
		"GenTime":  time.Now().Format(genTimeLayout),
		"IDLBytes": idlBytes,
	})
	if err != nil {
		println("Failed to execute plugin template:", err.Error())
		os.Exit(1)
	}
	code := buf.Bytes()
	code, err = format.Source(buf.Bytes())
	if err != nil {
		println("Failed to format generated code:", err.Error())
		os.Exit(1)
	}

	outputFile := getOutputFile(request.AST, lastService, request.OutputPath)
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

func encodeIDLToBytes(ast *parser.Thrift) (string, error) {
	buf, err := meta.Marshal(ast)
	if err != nil {
		return "", fmt.Errorf("cannot marshal parser.Thrift: %w", err)
	}
	gzBuf := bytes.NewBuffer(nil)
	gzw := gzip.NewWriter(gzBuf)
	_, err = gzw.Write(buf)
	if err != nil {
		return "", fmt.Errorf("cannot compress IDL bytes: %w", err)
	}
	err = gzw.Close()
	if err != nil {
		return "", fmt.Errorf("cannot compress IDL bytes: %w", err)
	}

	strb := strings.Builder{}
	for i, x := range gzBuf.Bytes() {
		if i > 0 && i%16 == 0 {
			strb.WriteByte('\n')
		}
		strb.WriteString(fmt.Sprintf("0x%02x,", x))
	}
	return strb.String(), nil
}

func getOutputFile(ast *parser.Thrift, svc *parser.Service, outputPath string) string {
	var namespace string
	for _, ns := range ast.Namespaces {
		if ns.Language == "go" {
			namespace = strings.ReplaceAll(ns.Name, ".", "/")
		}
	}
	svcName := strings.ToLower(svc.Name)
	filename := "plugin-reflect-desc.go"
	descFile := filepath.Join(outputPath, namespace, svcName, filename)
	return descFile
}
