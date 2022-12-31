package main

var implTpl = `
// Code generated by thrift-gen-kitex-reflect. DO NOT EDIT.
package {{ .PkgName }}

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	idl "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl"
)

var pluginReflectImpl = struct {
	genTime string

	once    sync.Once
	err     error
	version string
}{
	genTime: "{{ .GenTime }}",
}

func ServeReflectServiceRequest(ctx context.Context, req interface{}, resp interface{}) error {
	impl := &pluginReflectImpl
	impl.once.Do(func() {
		err := idl.CheckReflectReqAndRespType(req, resp)
		if err != nil {
			impl.err = err
			return
		}
		vcsRev := "unknown"
		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, setting := range buildInfo.Settings {
				if setting.Key == "vcs.revision" && setting.Value != "" {
					vcsRev = setting.Value
				}
			}
		}
		impl.version = fmt.Sprintf("%s/%s", impl.genTime, vcsRev)
	})
	if err := impl.err; err != nil {
		return err
	}

	reqPayload, err := idl.UnmarshalReflectServiceReqPayload(req.(interface {
		GetPayload() []byte
	}).GetPayload())
	if err != nil {
		return fmt.Errorf("cannot unmarshal ReflectServiceReqPayload: %w", err)
	}
	var reqVersionTime = ""
	if len(reqPayload.ExistingIDLVersion) > 14 {
		reqVersionTime = reqPayload.ExistingIDLVersion[:14]
	}

	respPayload := &idl.ReflectServiceRespPayload{
		Version: impl.version,
		IDL:     nil,
	}
	if reqVersionTime != impl.genTime {
		respPayload.IDL = pluginReflectIDLRaw
	}
	payloadBuf, err := idl.MarshalReflectServiceRespPayload(respPayload)
	if err != nil {
		return fmt.Errorf("cannot marshal ReflectServiceRespPayload: %w", err)
	}
	resp.(interface {
		SetPayload(val []byte)
	}).SetPayload(payloadBuf)
	return nil
}

func GetIDLGzipBytes() []byte {
	return pluginReflectIDLRaw
}

var pluginReflectIDLRaw = []byte{
{{ .IDLBytes }}
}
`
