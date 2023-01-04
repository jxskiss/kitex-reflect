namespace go kitexreflectidl

struct ReflectServiceReqPayload {
    1: string ExistingIDLVersion
}

struct ReflectServiceRespPayload {
    1: string Version
    2: bool IsCombineService

    15: binary IDL
}

struct ReflectServiceRequest {
    1: binary Payload; // ReflectServiceReqPayload
}

struct ReflectServiceResponse {
    1: binary Payload; // ReflectServiceRespPayload
}

service ReflectionService {
    ReflectServiceResponse ReflectService(1: ReflectServiceRequest request)
}
