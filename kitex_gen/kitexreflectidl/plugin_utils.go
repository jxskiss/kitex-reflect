package kitexreflectidl

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/apache/thrift/lib/go/thrift"
)

// MarshalReflectServiceReqPayload encodes a ReflectServiceReqPayload with binary protocol.
func MarshalReflectServiceReqPayload(payload *ReflectServiceReqPayload) ([]byte, error) {
	s := thrift.NewTSerializer()
	return s.Write(context.Background(), payload)
}

// UnmarshalReflectServiceReqPayload decodes a ReflectServiceReqPayload with binary protocol.
func UnmarshalReflectServiceReqPayload(bs []byte) (*ReflectServiceReqPayload, error) {
	payload := &ReflectServiceReqPayload{}
	d := thrift.NewTDeserializer()
	err := d.Read(payload, bs)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// MarshalReflectServiceRespPayload encodes a ReflectServiceRespPayload with binary protocol.
func MarshalReflectServiceRespPayload(payload *ReflectServiceRespPayload) ([]byte, error) {
	s := thrift.NewTSerializer()
	return s.Write(context.Background(), payload)
}

// UnmarshalReflectServiceRespPayload decodes a ReflectServiceRespPayload with binary protocol.
func UnmarshalReflectServiceRespPayload(bs []byte) (*ReflectServiceRespPayload, error) {
	payload := &ReflectServiceRespPayload{}
	d := thrift.NewTDeserializer()
	err := d.Read(payload, bs)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// CheckReflectReqAndRespType checks request and response type are compatible
// with the plugin's idl.
func CheckReflectReqAndRespType(req interface{}, resp interface{}) error {
	reqTyp := reflect.TypeOf(req)
	respTyp := reflect.TypeOf(resp)
	if (reqTyp.Kind() != reflect.Pointer && reqTyp.Elem().Kind() == reflect.Struct) ||
		(respTyp.Kind() != reflect.Pointer && respTyp.Elem().Kind() == reflect.Struct) {
		return errors.New("reflect request and response must be pointer to struct")
	}
	_, ok := req.(interface {
		GetPayload() (v []byte)
	})
	if !ok {
		return errors.New("reflect request does not have method GetPayload")
	}
	_, ok = req.(interface {
		SetPayload(val []byte)
	})
	if !ok {
		return errors.New("reflect response does not have method SetPayload")
	}

	byteSliceTyp := reflect.TypeOf([]byte(nil))
	getThriftFieldID := func(field reflect.StructField) string {
		tag := field.Tag.Get("thrift")
		parts := strings.Split(tag, ",")
		if len(parts) >= 2 {
			return parts[1] // the field's id number
		}
		return ""
	}

	reqPayloadField, ok := reqTyp.Elem().FieldByName("Payload")
	if !ok {
		return errors.New("reflect request does not have field Payload")
	}
	reqPayloadFieldID := getThriftFieldID(reqPayloadField)
	if reqPayloadField.Type != byteSliceTyp || reqPayloadFieldID != "1" {
		return errors.New("reflect request field Payload definition does not match")
	}

	respPayloadField, ok := respTyp.Elem().FieldByName("Payload")
	if !ok {
		return errors.New("reflect response does not have field Payload")
	}
	respPayloadFieldID := getThriftFieldID(respPayloadField)
	if respPayloadField.Type != byteSliceTyp || respPayloadFieldID != "1" {
		return errors.New("reflect response field Payload definition does not match")
	}

	return nil
}
