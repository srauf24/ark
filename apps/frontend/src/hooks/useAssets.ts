import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "@/api";
import type {
    Asset,
    CreateAssetRequest,
    UpdateAssetRequest,
    AssetQueryParams,
} from "@/types";

/**
 * Query keys factory for asset-related queries
 * Provides a centralized way to manage query keys for cache invalidation
 */
export const assetKeys = {
    all: ["assets"] as const,
    lists: () => [...assetKeys.all, "list"] as const,
    list: (params?: AssetQueryParams) => [...assetKeys.lists(), params] as const,
    details: () => [...assetKeys.all, "detail"] as const,
    detail: (id: string) => [...assetKeys.details(), id] as const,
};

/**
 * Hook to fetch a paginated list of assets
 * @param params - Query parameters for filtering, sorting, and pagination
 */
export function useAssets(params?: AssetQueryParams) {
    const apiClient = useApiClient();

    return useQuery({
        queryKey: assetKeys.list(params),
        queryFn: async () => {
            const response = await apiClient.Assets.listAssets({
                query: params || {},
            });

            if (response.status === 200) {
                return response.body;
            }

            throw new Error("Failed to fetch assets");
        },
        staleTime: 1000 * 60 * 5, // 5 minutes
    });
}

/**
 * Hook to fetch a single asset by ID
 * @param id - Asset ID
 */
export function useAsset(id: string) {
    const apiClient = useApiClient();

    return useQuery({
        queryKey: assetKeys.detail(id),
        queryFn: async () => {
            const response = await apiClient.Assets.getAssetById({
                params: { id },
            });

            if (response.status === 200) {
                return response.body;
            }

            throw new Error(`Failed to fetch asset ${id}`);
        },
        enabled: !!id,
        staleTime: 1000 * 60 * 5, // 5 minutes
    });
}

/**
 * Hook to create a new asset
 */
export function useCreateAsset() {
    const apiClient = useApiClient();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (data: CreateAssetRequest) => {
            const response = await apiClient.Assets.createAsset({
                body: data,
            });

            if (response.status === 201) {
                return response.body;
            }

            throw new Error("Failed to create asset");
        },
        onSuccess: () => {
            // Invalidate all asset list queries to refetch with new data
            queryClient.invalidateQueries({ queryKey: assetKeys.lists() });
        },
    });
}

/**
 * Hook to update an existing asset
 */
export function useUpdateAsset() {
    const apiClient = useApiClient();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ id, data }: { id: string; data: UpdateAssetRequest }) => {
            const response = await apiClient.Assets.updateAsset({
                params: { id },
                body: data,
            });

            if (response.status === 200) {
                return response.body;
            }

            throw new Error(`Failed to update asset ${id}`);
        },
        onSuccess: (asset: Asset) => {
            // Invalidate the specific asset detail query
            queryClient.invalidateQueries({ queryKey: assetKeys.detail(asset.id) });
            // Invalidate all asset list queries
            queryClient.invalidateQueries({ queryKey: assetKeys.lists() });
        },
    });
}

/**
 * Hook to delete an asset
 */
export function useDeleteAsset() {
    const apiClient = useApiClient();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (id: string) => {
            const response = await apiClient.Assets.deleteAsset({
                params: { id },
            });

            if (response.status === 204) {
                return id;
            }

            throw new Error(`Failed to delete asset ${id}`);
        },
        onSuccess: () => {
            // Invalidate all asset list queries to refetch without deleted asset
            queryClient.invalidateQueries({ queryKey: assetKeys.lists() });
        },
    });
}
