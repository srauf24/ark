import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "@/api";
import { queryKeys } from "@/api/utils";
import { toast } from "sonner";

export const useLogs = (assetId: string, params?: { limit?: number; offset?: number; tags?: string[]; search?: string }) => {
    const apiClient = useApiClient();

    return useQuery({
        queryKey: queryKeys.logs.list(assetId, params),
        queryFn: async () => {
            const response = await apiClient.Logs.listLogsByAsset({
                params: { id: assetId },
                query: params || {},
            });

            if (response.status !== 200) {
                throw new Error("Failed to fetch logs");
            }

            return response.body;
        },
        enabled: !!assetId,
    });
};

export const useCreateLog = (assetId: string) => {
    const apiClient = useApiClient();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (data: { content: string; tags?: string[] }) => {
            const response = await apiClient.Logs.createLog({
                params: { id: assetId },
                body: data,
            });

            if (response.status !== 201) {
                throw new Error("Failed to create log");
            }

            return response.body;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: queryKeys.logs.all(assetId) });
            toast.success("Log created successfully");
        },
        onError: () => {
            toast.error("Failed to create log");
        },
    });
};

export const useUpdateLog = (assetId: string) => {
    const apiClient = useApiClient();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ id, data }: { id: string; data: { content?: string; tags?: string[] } }) => {
            const response = await apiClient.Logs.updateLog({
                params: { id },
                body: data,
            });

            if (response.status !== 200) {
                throw new Error("Failed to update log");
            }

            return response.body;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: queryKeys.logs.all(assetId) });
            toast.success("Log updated successfully");
        },
        onError: () => {
            toast.error("Failed to update log");
        },
    });
};

export const useDeleteLog = (assetId: string) => {
    const apiClient = useApiClient();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (id: string) => {
            const response = await apiClient.Logs.deleteLog({
                params: { id },
            });

            if (response.status !== 204) {
                throw new Error("Failed to delete log");
            }

            return true;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: queryKeys.logs.all(assetId) });
            toast.success("Log deleted successfully");
        },
        onError: () => {
            toast.error("Failed to delete log");
        },
    });
};
