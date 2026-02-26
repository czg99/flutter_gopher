interface FgRequest {
  method: number;
  data: ArrayBuffer;
}

interface FgResponse {
  data: ArrayBuffer;
  error: ArrayBuffer;
}

export const callGoMethod: (request: FgRequest) => FgResponse;
export const callDartMethod: (request: FgRequest) => FgResponse;
