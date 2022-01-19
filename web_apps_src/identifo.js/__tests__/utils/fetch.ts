export function setupFetchStub<T>(data: any) {
  return function fetchStub(input: RequestInfo, init?: RequestInit | undefined): Promise<Response> {
    return new Promise((resolve) => {
      resolve({
        ok: true,
        json: () => Promise.resolve({ ...data }),
      } as any);
    });
  };
}
