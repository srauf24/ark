import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "@/api";
import type {
    AssetLog,
    CreateLogRequest,
    UpdateLogRequest,
    LogQueryParams,
} from "@/types";

/**
 * Query keys factory for log-related queries
 */
export const logKeys = {
    all: ["logs"] as const,
    lists: () => [...logKeys.all, "list"] as const,
    list: (assetId: string, params?: LogQueryParams) =>
        [...logKeys.lists(), assetId, params] as const,
    details: () => [...logKeys.all, "detail"] as const,
    detail: (id: string) => [...logKeys.details(), id] as const,
};

/**
 * Hook to fetch a paginated list of logs for an asset
 * @param assetId - The ID of the asset to fetch logs for
 * @param params - Query parameters
 */
export function useLogs(assetId: string, params?: LogQueryParams) {
    const apiClient = useApiClient();

    return useQuery({
        queryKey: logKeys.list(assetId, params),
        queryFn: async () => {
            const response = await apiClient.Logs.listLogsByAsset({
                params: { id: assetId },
                query: params || {},
            });

            if (response.status === 200) {
                return response.body;
            }

            throw new Error("Failed to fetch logs");
        },
        enabled: !!assetId,
        staleTime: 1000 * 60 * 2, // 2 minutes
    });
}

/**
 * Hook to fetch a single log by ID
 * @param id - Log ID
 */
export function useLog(id: string) {
    const apiClient = useApiClient();

    return useQuery({
        queryKey: logKeys.detail(id),
        queryFn: async () => {
            const response = await apiClient.Logs.getLogById({
                params: { id },
            });

            if (response.status === 200) {
                return response.body.data;
            }

            throw new Error(`Failed to fetch log ${id}`);
        },
        enabled: !!id,
        staleTime: 1000 * 60 * 5, // 5 minutes
    });
}

/**
 * Hook to create a new log
 */
export function useCreateLog() {
    const apiClient = useApiClient();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ assetId, data }: { assetId: string; data: CreateLogRequest }) => {
            const response = await apiClient.Logs.createLog({
                params: { id: assetId },
                body: data,
            });

            if (response.status === 201) {
                return { log: response.body.data, assetId };
            }

            throw new Error("Failed to create log");
        },
        onSuccess: ({ assetId }) => {
            // Invalidate logs list for this specific asset
            queryClient.invalidateQueries({
                queryKey: [...logKeys.lists(), assetId]
            });
        },
    });
}

/**
 * Hook to update an existing log
 */
export function useUpdateLog() {
    const apiClient = useApiClient();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ id, data }: { id: string; data: UpdateLogRequest }) => {
            const response = await apiClient.Logs.updateLog({
                params: { id },
                body: data,
            });

            if (response.status === 200) {
                return response.body.data;
            }

            throw new Error(`Failed to update log ${id}`);
        },
        onSuccess: (log: AssetLog) => {
            // Invalidate specific log detail
            queryClient.invalidateQueries({ queryKey: logKeys.detail(log.id) });
            // Invalidate logs list for the asset this log belongs to
            queryClient.invalidateQueries({
                queryKey: [...logKeys.lists(), log.asset_id]
            });
        },
    });
}

/**
 * Hook to delete a log
 */
export function useDeleteLog() {
    const apiClient = useApiClient();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ id, assetId }: { id: string; assetId: string }) => {
            const response = await apiClient.Logs.deleteLog({
                params: { id },
            });

            if (response.status === 204) {
                return { id, assetId };
            }

            throw new Error(`Failed to delete log ${id}`);
        },
        onSuccess: ({ assetId }) => {
            // Invalidate logs list for the asset
            queryClient.invalidateQueries({
                queryKey: [...logKeys.lists(), assetId]
            });
        },
    });
}
