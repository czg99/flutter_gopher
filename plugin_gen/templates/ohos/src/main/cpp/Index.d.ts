interface FgRequest {
  method: number;
  data: ArrayBuffer | undefined;
}

interface FgResponse {
  data: ArrayBuffer | undefined;
  error: ArrayBuffer | undefined;
}

export const callGoMethod: (request: FgRequest) => FgResponse;
export const callGoMethodAsync: (request: FgRequest) => Promise<FgResponse>;
export const callDartMethod: (request: FgRequest) => void;
export const initPlatformMethodHandle:(callback: (request: FgRequest) => Promise<FgResponse>) => void;
