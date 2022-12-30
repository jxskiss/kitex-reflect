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
    1: optional string Name
    2: optional string Alias
    3: optional i32 ID
    4: optional bool Required
    5: optional bool Optional
    6: optional bool IsException
    7: optional string DefaultValue
    8: optional TypeDesc Type
    9: optional list<Annotation> Annotations
}

struct TypeDesc {
    1: optional string Name
    2: optional Type Type
    3: optional TypeDesc Key       // for map key
    4: optional TypeDesc Elem      // for slice or map element
    5: optional StructDesc Struct  // for struct
    6: optional i32 StructIdx      // for struct
    7: optional bool IsRequestBase
}

struct StructDesc {
    1: optional string Name
    2: optional list<FieldDesc> Fields
    3: optional list<Annotation> Annotations
    4: optional string UniqueKey
}

struct FunctionDesc {
    1: optional string Name
    2: optional bool Oneway
    3: optional bool HasRequestBase
    4: optional TypeDesc Request
    5: optional TypeDesc Response
    6: optional list<Annotation> Annotations
}

struct ServiceDesc {
    1: string Name
    2: list<FunctionDesc> Functions
    3: list<StructDesc> StructList
}

struct ReflectServiceReqPayload {
}

struct ReflectServiceRespPayload {
    1: string Version
    2: ServiceDesc ServiceDesc
}

struct ReflectServiceRequest {
    1: binary Payload; // ReflectServiceReqPayload
}

struct ReflectServiceResponse {
    1: binary Payload; // ReflectServiceRespPayload
}

service ReflectionService {
    ReflectServiceResponse ReflectService(1: ReflectServiceRequest req)
}
