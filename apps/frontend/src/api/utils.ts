export const queryKeys = {
  assets: {
    all: ["assets"] as const,
    list: (params: any) => ["assets", "list", params] as const,
    detail: (id: string) => ["assets", "detail", id] as const,
  },
  logs: {
    all: (assetId: string) => ["logs", assetId] as const,
    list: (assetId: string, params: any) => ["logs", assetId, "list", params] as const,
  },
};

