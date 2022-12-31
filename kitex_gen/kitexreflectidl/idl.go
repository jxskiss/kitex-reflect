// Code generated by thriftgo (0.2.4). DO NOT EDIT.

package kitexreflectidl

import (
	"bytes"
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"strings"
)

type ReflectServiceReqPayload struct {
	ExistingIDLVersion string `thrift:"ExistingIDLVersion,1" frugal:"1,default,string" json:"ExistingIDLVersion"`
}

func NewReflectServiceReqPayload() *ReflectServiceReqPayload {
	return &ReflectServiceReqPayload{}
}

func (p *ReflectServiceReqPayload) InitDefault() {
	*p = ReflectServiceReqPayload{}
}

func (p *ReflectServiceReqPayload) GetExistingIDLVersion() (v string) {
	return p.ExistingIDLVersion
}
func (p *ReflectServiceReqPayload) SetExistingIDLVersion(val string) {
	p.ExistingIDLVersion = val
}

var fieldIDToName_ReflectServiceReqPayload = map[int16]string{
	1: "ExistingIDLVersion",
}

func (p *ReflectServiceReqPayload) Read(iprot thrift.TProtocol) (err error) {

	var fieldTypeId thrift.TType
	var fieldId int16

	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}

	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}

		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}

	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_ReflectServiceReqPayload[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)

ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
}

func (p *ReflectServiceReqPayload) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		p.ExistingIDLVersion = v
	}
	return nil
}

func (p *ReflectServiceReqPayload) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("ReflectServiceReqPayload"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}

	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *ReflectServiceReqPayload) writeField1(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("ExistingIDLVersion", thrift.STRING, 1); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.ExistingIDLVersion); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *ReflectServiceReqPayload) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ReflectServiceReqPayload(%+v)", *p)
}

func (p *ReflectServiceReqPayload) DeepEqual(ano *ReflectServiceReqPayload) bool {
	if p == ano {
		return true
	} else if p == nil || ano == nil {
		return false
	}
	if !p.Field1DeepEqual(ano.ExistingIDLVersion) {
		return false
	}
	return true
}

func (p *ReflectServiceReqPayload) Field1DeepEqual(src string) bool {

	if strings.Compare(p.ExistingIDLVersion, src) != 0 {
		return false
	}
	return true
}

type ReflectServiceRespPayload struct {
	Version string `thrift:"Version,1" frugal:"1,default,string" json:"Version"`
	IDL     []byte `thrift:"IDL,2" frugal:"2,default,binary" json:"IDL"`
}

func NewReflectServiceRespPayload() *ReflectServiceRespPayload {
	return &ReflectServiceRespPayload{}
}

func (p *ReflectServiceRespPayload) InitDefault() {
	*p = ReflectServiceRespPayload{}
}

func (p *ReflectServiceRespPayload) GetVersion() (v string) {
	return p.Version
}

func (p *ReflectServiceRespPayload) GetIDL() (v []byte) {
	return p.IDL
}
func (p *ReflectServiceRespPayload) SetVersion(val string) {
	p.Version = val
}
func (p *ReflectServiceRespPayload) SetIDL(val []byte) {
	p.IDL = val
}

var fieldIDToName_ReflectServiceRespPayload = map[int16]string{
	1: "Version",
	2: "IDL",
}

func (p *ReflectServiceRespPayload) Read(iprot thrift.TProtocol) (err error) {

	var fieldTypeId thrift.TType
	var fieldId int16

	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}

	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		case 2:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField2(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}

		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}

	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_ReflectServiceRespPayload[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)

ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
}

func (p *ReflectServiceRespPayload) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		p.Version = v
	}
	return nil
}

func (p *ReflectServiceRespPayload) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return err
	} else {
		p.IDL = []byte(v)
	}
	return nil
}

func (p *ReflectServiceRespPayload) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("ReflectServiceRespPayload"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}
		if err = p.writeField2(oprot); err != nil {
			fieldId = 2
			goto WriteFieldError
		}

	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *ReflectServiceRespPayload) writeField1(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("Version", thrift.STRING, 1); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.Version); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *ReflectServiceRespPayload) writeField2(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("IDL", thrift.STRING, 2); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteBinary([]byte(p.IDL)); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 end error: ", p), err)
}

func (p *ReflectServiceRespPayload) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ReflectServiceRespPayload(%+v)", *p)
}

func (p *ReflectServiceRespPayload) DeepEqual(ano *ReflectServiceRespPayload) bool {
	if p == ano {
		return true
	} else if p == nil || ano == nil {
		return false
	}
	if !p.Field1DeepEqual(ano.Version) {
		return false
	}
	if !p.Field2DeepEqual(ano.IDL) {
		return false
	}
	return true
}

func (p *ReflectServiceRespPayload) Field1DeepEqual(src string) bool {

	if strings.Compare(p.Version, src) != 0 {
		return false
	}
	return true
}
func (p *ReflectServiceRespPayload) Field2DeepEqual(src []byte) bool {

	if bytes.Compare(p.IDL, src) != 0 {
		return false
	}
	return true
}

type ReflectServiceRequest struct {
	Payload []byte `thrift:"Payload,1" frugal:"1,default,binary" json:"Payload"`
}

func NewReflectServiceRequest() *ReflectServiceRequest {
	return &ReflectServiceRequest{}
}

func (p *ReflectServiceRequest) InitDefault() {
	*p = ReflectServiceRequest{}
}

func (p *ReflectServiceRequest) GetPayload() (v []byte) {
	return p.Payload
}
func (p *ReflectServiceRequest) SetPayload(val []byte) {
	p.Payload = val
}

var fieldIDToName_ReflectServiceRequest = map[int16]string{
	1: "Payload",
}

func (p *ReflectServiceRequest) Read(iprot thrift.TProtocol) (err error) {

	var fieldTypeId thrift.TType
	var fieldId int16

	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}

	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}

		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}

	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_ReflectServiceRequest[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)

ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
}

func (p *ReflectServiceRequest) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return err
	} else {
		p.Payload = []byte(v)
	}
	return nil
}

func (p *ReflectServiceRequest) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("ReflectServiceRequest"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}

	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *ReflectServiceRequest) writeField1(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("Payload", thrift.STRING, 1); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteBinary([]byte(p.Payload)); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *ReflectServiceRequest) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ReflectServiceRequest(%+v)", *p)
}

func (p *ReflectServiceRequest) DeepEqual(ano *ReflectServiceRequest) bool {
	if p == ano {
		return true
	} else if p == nil || ano == nil {
		return false
	}
	if !p.Field1DeepEqual(ano.Payload) {
		return false
	}
	return true
}

func (p *ReflectServiceRequest) Field1DeepEqual(src []byte) bool {

	if bytes.Compare(p.Payload, src) != 0 {
		return false
	}
	return true
}

type ReflectServiceResponse struct {
	Payload []byte `thrift:"Payload,1" frugal:"1,default,binary" json:"Payload"`
}

func NewReflectServiceResponse() *ReflectServiceResponse {
	return &ReflectServiceResponse{}
}

func (p *ReflectServiceResponse) InitDefault() {
	*p = ReflectServiceResponse{}
}

func (p *ReflectServiceResponse) GetPayload() (v []byte) {
	return p.Payload
}
func (p *ReflectServiceResponse) SetPayload(val []byte) {
	p.Payload = val
}

var fieldIDToName_ReflectServiceResponse = map[int16]string{
	1: "Payload",
}

func (p *ReflectServiceResponse) Read(iprot thrift.TProtocol) (err error) {

	var fieldTypeId thrift.TType
	var fieldId int16

	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}

	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}

		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}

	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_ReflectServiceResponse[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)

ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
}

func (p *ReflectServiceResponse) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return err
	} else {
		p.Payload = []byte(v)
	}
	return nil
}

func (p *ReflectServiceResponse) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("ReflectServiceResponse"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}

	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *ReflectServiceResponse) writeField1(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("Payload", thrift.STRING, 1); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteBinary([]byte(p.Payload)); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *ReflectServiceResponse) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ReflectServiceResponse(%+v)", *p)
}

func (p *ReflectServiceResponse) DeepEqual(ano *ReflectServiceResponse) bool {
	if p == ano {
		return true
	} else if p == nil || ano == nil {
		return false
	}
	if !p.Field1DeepEqual(ano.Payload) {
		return false
	}
	return true
}

func (p *ReflectServiceResponse) Field1DeepEqual(src []byte) bool {

	if bytes.Compare(p.Payload, src) != 0 {
		return false
	}
	return true
}

type ReflectionService interface {
	ReflectService(ctx context.Context, request *ReflectServiceRequest) (r *ReflectServiceResponse, err error)
}

type ReflectionServiceClient struct {
	c thrift.TClient
}

func NewReflectionServiceClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *ReflectionServiceClient {
	return &ReflectionServiceClient{
		c: thrift.NewTStandardClient(f.GetProtocol(t), f.GetProtocol(t)),
	}
}

func NewReflectionServiceClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *ReflectionServiceClient {
	return &ReflectionServiceClient{
		c: thrift.NewTStandardClient(iprot, oprot),
	}
}

func NewReflectionServiceClient(c thrift.TClient) *ReflectionServiceClient {
	return &ReflectionServiceClient{
		c: c,
	}
}

func (p *ReflectionServiceClient) Client_() thrift.TClient {
	return p.c
}

func (p *ReflectionServiceClient) ReflectService(ctx context.Context, request *ReflectServiceRequest) (r *ReflectServiceResponse, err error) {
	var _args ReflectionServiceReflectServiceArgs
	_args.Request = request
	var _result ReflectionServiceReflectServiceResult
	if err = p.Client_().Call(ctx, "ReflectService", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

type ReflectionServiceProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      ReflectionService
}

func (p *ReflectionServiceProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *ReflectionServiceProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *ReflectionServiceProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewReflectionServiceProcessor(handler ReflectionService) *ReflectionServiceProcessor {
	self := &ReflectionServiceProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self.AddToProcessorMap("ReflectService", &reflectionServiceProcessorReflectService{handler: handler})
	return self
}
func (p *ReflectionServiceProcessor) Process(ctx context.Context, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return false, err
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(ctx, seqId, iprot, oprot)
	}
	iprot.Skip(thrift.STRUCT)
	iprot.ReadMessageEnd()
	x := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	x.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush(ctx)
	return false, x
}

type reflectionServiceProcessorReflectService struct {
	handler ReflectionService
}

func (p *reflectionServiceProcessorReflectService) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := ReflectionServiceReflectServiceArgs{}
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("ReflectService", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush(ctx)
		return false, err
	}

	iprot.ReadMessageEnd()
	var err2 error
	result := ReflectionServiceReflectServiceResult{}
	var retval *ReflectServiceResponse
	if retval, err2 = p.handler.ReflectService(ctx, args.Request); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing ReflectService: "+err2.Error())
		oprot.WriteMessageBegin("ReflectService", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush(ctx)
		return true, err2
	} else {
		result.Success = retval
	}
	if err2 = oprot.WriteMessageBegin("ReflectService", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(ctx); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

type ReflectionServiceReflectServiceArgs struct {
	Request *ReflectServiceRequest `thrift:"request,1" frugal:"1,default,ReflectServiceRequest" json:"request"`
}

func NewReflectionServiceReflectServiceArgs() *ReflectionServiceReflectServiceArgs {
	return &ReflectionServiceReflectServiceArgs{}
}

func (p *ReflectionServiceReflectServiceArgs) InitDefault() {
	*p = ReflectionServiceReflectServiceArgs{}
}

var ReflectionServiceReflectServiceArgs_Request_DEFAULT *ReflectServiceRequest

func (p *ReflectionServiceReflectServiceArgs) GetRequest() (v *ReflectServiceRequest) {
	if !p.IsSetRequest() {
		return ReflectionServiceReflectServiceArgs_Request_DEFAULT
	}
	return p.Request
}
func (p *ReflectionServiceReflectServiceArgs) SetRequest(val *ReflectServiceRequest) {
	p.Request = val
}

var fieldIDToName_ReflectionServiceReflectServiceArgs = map[int16]string{
	1: "request",
}

func (p *ReflectionServiceReflectServiceArgs) IsSetRequest() bool {
	return p.Request != nil
}

func (p *ReflectionServiceReflectServiceArgs) Read(iprot thrift.TProtocol) (err error) {

	var fieldTypeId thrift.TType
	var fieldId int16

	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}

	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRUCT {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}

		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}

	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_ReflectionServiceReflectServiceArgs[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)

ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
}

func (p *ReflectionServiceReflectServiceArgs) ReadField1(iprot thrift.TProtocol) error {
	p.Request = NewReflectServiceRequest()
	if err := p.Request.Read(iprot); err != nil {
		return err
	}
	return nil
}

func (p *ReflectionServiceReflectServiceArgs) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("ReflectService_args"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}

	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *ReflectionServiceReflectServiceArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("request", thrift.STRUCT, 1); err != nil {
		goto WriteFieldBeginError
	}
	if err := p.Request.Write(oprot); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *ReflectionServiceReflectServiceArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ReflectionServiceReflectServiceArgs(%+v)", *p)
}

func (p *ReflectionServiceReflectServiceArgs) DeepEqual(ano *ReflectionServiceReflectServiceArgs) bool {
	if p == ano {
		return true
	} else if p == nil || ano == nil {
		return false
	}
	if !p.Field1DeepEqual(ano.Request) {
		return false
	}
	return true
}

func (p *ReflectionServiceReflectServiceArgs) Field1DeepEqual(src *ReflectServiceRequest) bool {

	if !p.Request.DeepEqual(src) {
		return false
	}
	return true
}

type ReflectionServiceReflectServiceResult struct {
	Success *ReflectServiceResponse `thrift:"success,0,optional" frugal:"0,optional,ReflectServiceResponse" json:"success,omitempty"`
}

func NewReflectionServiceReflectServiceResult() *ReflectionServiceReflectServiceResult {
	return &ReflectionServiceReflectServiceResult{}
}

func (p *ReflectionServiceReflectServiceResult) InitDefault() {
	*p = ReflectionServiceReflectServiceResult{}
}

var ReflectionServiceReflectServiceResult_Success_DEFAULT *ReflectServiceResponse

func (p *ReflectionServiceReflectServiceResult) GetSuccess() (v *ReflectServiceResponse) {
	if !p.IsSetSuccess() {
		return ReflectionServiceReflectServiceResult_Success_DEFAULT
	}
	return p.Success
}
func (p *ReflectionServiceReflectServiceResult) SetSuccess(x interface{}) {
	p.Success = x.(*ReflectServiceResponse)
}

var fieldIDToName_ReflectionServiceReflectServiceResult = map[int16]string{
	0: "success",
}

func (p *ReflectionServiceReflectServiceResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *ReflectionServiceReflectServiceResult) Read(iprot thrift.TProtocol) (err error) {

	var fieldTypeId thrift.TType
	var fieldId int16

	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}

	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 0:
			if fieldTypeId == thrift.STRUCT {
				if err = p.ReadField0(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}

		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}

	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_ReflectionServiceReflectServiceResult[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)

ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
}

func (p *ReflectionServiceReflectServiceResult) ReadField0(iprot thrift.TProtocol) error {
	p.Success = NewReflectServiceResponse()
	if err := p.Success.Read(iprot); err != nil {
		return err
	}
	return nil
}

func (p *ReflectionServiceReflectServiceResult) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("ReflectService_result"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField0(oprot); err != nil {
			fieldId = 0
			goto WriteFieldError
		}

	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *ReflectionServiceReflectServiceResult) writeField0(oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err = oprot.WriteFieldBegin("success", thrift.STRUCT, 0); err != nil {
			goto WriteFieldBeginError
		}
		if err := p.Success.Write(oprot); err != nil {
			return err
		}
		if err = oprot.WriteFieldEnd(); err != nil {
			goto WriteFieldEndError
		}
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 0 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 0 end error: ", p), err)
}

func (p *ReflectionServiceReflectServiceResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ReflectionServiceReflectServiceResult(%+v)", *p)
}

func (p *ReflectionServiceReflectServiceResult) DeepEqual(ano *ReflectionServiceReflectServiceResult) bool {
	if p == ano {
		return true
	} else if p == nil || ano == nil {
		return false
	}
	if !p.Field0DeepEqual(ano.Success) {
		return false
	}
	return true
}

func (p *ReflectionServiceReflectServiceResult) Field0DeepEqual(src *ReflectServiceResponse) bool {

	if !p.Success.DeepEqual(src) {
		return false
	}
	return true
}
