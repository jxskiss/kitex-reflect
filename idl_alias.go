package kitexreflect

import (
	idl "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl"
)

type (
	ReflectServiceReqPayload  = idl.ReflectServiceReqPayload
	ReflectServiceRespPayload = idl.ReflectServiceRespPayload
	ReflectServiceRequest     = idl.ReflectServiceRequest
	ReflectServiceResponse    = idl.ReflectServiceResponse
)

var (
	CheckReflectReqAndRespType = idl.CheckReflectReqAndRespType

	MarshalReflectServiceReqPayload   = idl.MarshalReflectServiceReqPayload
	UnmarshalReflectServiceReqPayload = idl.UnmarshalReflectServiceReqPayload

	MarshalReflectServiceRespPayload   = idl.MarshalReflectServiceRespPayload
	UnmarshalReflectServiceRespPayload = idl.UnmarshalReflectServiceRespPayload
)
