import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import {
    useAssets,
    useAsset,
    useCreateAsset,
    useUpdateAsset,
    useDeleteAsset,
    assetKeys,
} from "./useAssets";

// Mock the API client
const mockApiClient = {
    Assets: {
        listAssets: vi.fn(),
        getAssetById: vi.fn(),
        createAsset: vi.fn(),
        updateAsset: vi.fn(),
        deleteAsset: vi.fn(),
    },
};

vi.mock("@/api", () => ({
    useApiClient: () => mockApiClient,
}));

// Test wrapper with QueryClient
function createWrapper() {
    const queryClient = new QueryClient({
        defaultOptions: {
            queries: { retry: false },
            mutations: { retry: false },
        },
    });

    return function Wrapper({ children }: { children: ReactNode }) {
    return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
    };
}

describe("useAssets hooks", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    describe("assetKeys", () => {
        it("should generate correct query keys", () => {
            expect(assetKeys.all).toEqual(["assets"]);
            expect(assetKeys.lists()).toEqual(["assets", "list"]);
            expect(assetKeys.list({ limit: 10 })).toEqual(["assets", "list", { limit: 10 }]);
            expect(assetKeys.details()).toEqual(["assets", "detail"]);
            expect(assetKeys.detail("123")).toEqual(["assets", "detail", "123"]);
        });
    });

    describe("useAssets", () => {
        it("should fetch assets list successfully", async () => {
            const mockAssets = {
                data: [
                    {
                        id: "1",
                        user_id: "user1",
                        name: "Server 1",
                        type: "server" as const,
                        created_at: "2024-01-01T00:00:00Z",
                        updated_at: "2024-01-01T00:00:00Z",
                    },
                ],
                total: 1,
                limit: 50,
                offset: 0,
            };

            mockApiClient.Assets.listAssets.mockResolvedValue({
                status: 200,
                body: mockAssets,
            });

            const { result } = renderHook(() => useAssets(), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));

            expect(result.current.data).toEqual(mockAssets);
            expect(mockApiClient.Assets.listAssets).toHaveBeenCalledWith({
                query: {},
            });
        });

        it("should fetch assets with query parameters", async () => {
            const params = { limit: 10, offset: 0, type: "server" as const };
            mockApiClient.Assets.listAssets.mockResolvedValue({
                status: 200,
                body: { data: [], total: 0, limit: 10, offset: 0 },
            });

            const { result } = renderHook(() => useAssets(params), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));

            expect(mockApiClient.Assets.listAssets).toHaveBeenCalledWith({
                query: params,
            });
        });

        it("should handle fetch error", async () => {
            mockApiClient.Assets.listAssets.mockResolvedValue({
                status: 500,
                body: { error: "Server error" },
            });

            const { result } = renderHook(() => useAssets(), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isError).toBe(true));
            expect(result.current.error).toBeDefined();
        });
    });

    describe("useAsset", () => {
        it("should fetch single asset successfully", async () => {
            const mockAsset = {
                id: "123",
                user_id: "user1",
                name: "Test Server",
                type: "server" as const,
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-01T00:00:00Z",
            };

            mockApiClient.Assets.getAssetById.mockResolvedValue({
                status: 200,
                body: { data: mockAsset },
            });

            const { result } = renderHook(() => useAsset("123"), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));

            expect(result.current.data).toEqual(mockAsset);
            expect(mockApiClient.Assets.getAssetById).toHaveBeenCalledWith({
                params: { id: "123" },
            });
        });

        it("should be disabled when no ID provided", () => {
            const { result } = renderHook(() => useAsset(""), {
                wrapper: createWrapper(),
            });

            expect(result.current.fetchStatus).toBe("idle");
            expect(mockApiClient.Assets.getAssetById).not.toHaveBeenCalled();
        });
    });

    describe("useCreateAsset", () => {
        it("should create asset successfully", async () => {
            const newAsset = {
                id: "new-123",
                user_id: "user1",
                name: "New Server",
                type: "server" as const,
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-01T00:00:00Z",
            };

            mockApiClient.Assets.createAsset.mockResolvedValue({
                status: 201,
                body: { data: newAsset },
            });

            const { result } = renderHook(() => useCreateAsset(), {
                wrapper: createWrapper(),
            });

            result.current.mutate({ name: "New Server", type: "server" });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));

            expect(result.current.data).toEqual(newAsset);
            expect(mockApiClient.Assets.createAsset).toHaveBeenCalledWith({
                body: { name: "New Server", type: "server" },
            });
        });

        it("should handle create error", async () => {
            mockApiClient.Assets.createAsset.mockResolvedValue({
                status: 400,
                body: { error: "Invalid data" },
            });

            const { result } = renderHook(() => useCreateAsset(), {
                wrapper: createWrapper(),
            });

            result.current.mutate({ name: "New Server" });

            await waitFor(() => expect(result.current.isError).toBe(true));
            expect(result.current.error).toBeDefined();
        });
    });

    describe("useUpdateAsset", () => {
        it("should update asset successfully", async () => {
            const updatedAsset = {
                id: "123",
                user_id: "user1",
                name: "Updated Server",
                type: "server" as const,
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-02T00:00:00Z",
            };

            mockApiClient.Assets.updateAsset.mockResolvedValue({
                status: 200,
                body: { data: updatedAsset },
            });

            const { result } = renderHook(() => useUpdateAsset(), {
                wrapper: createWrapper(),
            });

            result.current.mutate({ id: "123", data: { name: "Updated Server" } });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));

            expect(result.current.data).toEqual(updatedAsset);
            expect(mockApiClient.Assets.updateAsset).toHaveBeenCalledWith({
                params: { id: "123" },
                body: { name: "Updated Server" },
            });
        });
    });

    describe("useDeleteAsset", () => {
        it("should delete asset successfully", async () => {
            mockApiClient.Assets.deleteAsset.mockResolvedValue({
                status: 204,
                body: undefined,
            });

            const { result } = renderHook(() => useDeleteAsset(), {
                wrapper: createWrapper(),
            });

            result.current.mutate("123");

            await waitFor(() => expect(result.current.isSuccess).toBe(true));

            expect(result.current.data).toBe("123");
            expect(mockApiClient.Assets.deleteAsset).toHaveBeenCalledWith({
                params: { id: "123" },
            });
        });

        it("should handle delete error", async () => {
            mockApiClient.Assets.deleteAsset.mockResolvedValue({
                status: 404,
                body: { error: "Not found" },
            });

            const { result } = renderHook(() => useDeleteAsset(), {
                wrapper: createWrapper(),
            });

            result.current.mutate("123");

            await waitFor(() => expect(result.current.isError).toBe(true));
            expect(result.current.error).toBeDefined();
        });
    });
});
