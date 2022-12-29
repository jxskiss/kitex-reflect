namespace go kitexreflectidl

enum Type {
    STOP   = 0
    VOID   = 1
    BOOL   = 2
    BYTE   = 3 // BYTE, I08
    DOUBLE = 4
    // 5 unused
    I16    = 6
    I32    = 8
    I64    = 10
    STRING = 11 // STRING, UTF7
    STRUCT = 12
    MAP    = 13
    SET    = 14
    LIST   = 15
    UTF8   = 16
    UTF16  = 17
    // BINARY = 18   wrong and unusued
    JSON   = 19
}

struct Annotation {
    1: string Key
    2: list<string> Values
}

struct FieldDesc {
    1: string Name
    2: string Alias
    3: i32 ID
    4: bool Required
    5: bool Optional
    6: bool IsException
    7: optional string DefaultValue
    8: optional i32 TypeDescIndex
    9: optional TypeDesc Type
    10: list<Annotation> Annotations
}

struct TypeDesc {
    1: string Name
    2: Type Type
    3: optional TypeDesc Key       // for map key
    4: optional TypeDesc Elem      // for slice or map element
    5: optional StructDesc Struct  // for struct
    6: optional bool IsRequestBase
}

struct StructDesc {
    1: string Name
    2: list<FieldDesc> Fields
    3: list<Annotation> Annotations
}

struct FunctionDesc {
    1: string Name
    2: bool Oneway
    3: bool HasRequestBase
    4: TypeDesc Request
    5: TypeDesc Response
    6: list<Annotation> Annotations
}

struct ServiceDesc {
    1: string Name
    2: list<FunctionDesc> Functions
    3: list<TypeDesc> TypeDescList
}

struct ReflectServiceReqPayload {
}

struct ReflectServiceRespPayload {
    1: string Version
    2: ServiceDesc ServiceDesc
}

struct ReflectServiceRequest {
    1: binary Payload; // ReflectServiceReqData
}

struct ReflectServiceResponse {
    1: binary Payload; // ReflectServiceRespData
}

service ReflectionService {
    ReflectServiceResponse ReflectService(1: ReflectServiceRequest req)
}
