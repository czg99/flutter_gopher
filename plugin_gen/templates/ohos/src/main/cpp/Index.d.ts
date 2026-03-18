interface FgRequest {
  method: number;
  data: ArrayBuffer;
}

interface FgResponse {
  data: ArrayBuffer;
  error: ArrayBuffer;
}

export const callGoMethod: (request: FgRequest) => FgResponse;
export const callGoMethodAsync: (request: FgRequest) => Promise<FgResponse>;
export const callDartMethod: (request: FgRequest) => void;
export const initPlatformMethodHandle:(callback: (request: FgRequest) => FgResponse) => void;
