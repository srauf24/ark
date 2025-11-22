import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import {
    useLogs,
    useLog,
    useCreateLog,
    useUpdateLog,
    useDeleteLog,
    logKeys,
} from "./useLogs";

// Mock the API client
const mockApiClient = {
    Logs: {
        listLogsByAsset: vi.fn(),
        getLogById: vi.fn(),
        createLog: vi.fn(),
        updateLog: vi.fn(),
        deleteLog: vi.fn(),
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
        return (
            <QueryClientProvider client={queryClient}>
                {children}
            </QueryClientProvider>
        );
    };
}

describe("useLogs hooks", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    describe("logKeys", () => {
        it("should generate correct query keys", () => {
            expect(logKeys.all).toEqual(["logs"]);
            expect(logKeys.lists()).toEqual(["logs", "list"]);
            expect(logKeys.list("asset-1", { limit: 10 })).toEqual(["logs", "list", "asset-1", { limit: 10 }]);
            expect(logKeys.details()).toEqual(["logs", "detail"]);
            expect(logKeys.detail("log-1")).toEqual(["logs", "detail", "log-1"]);
        });
    });

    describe("useLogs", () => {
        it("should fetch logs list successfully", async () => {
            const mockLogs = {
                data: [
                    {
                        id: "log-1",
                        asset_id: "asset-1",
                        user_id: "user1",
                        content: "Log content",
                        created_at: "2024-01-01T00:00:00Z",
                        updated_at: "2024-01-01T00:00:00Z",
                    },
                ],
                total: 1,
                limit: 50,
                offset: 0,
            };

            mockApiClient.Logs.listLogsByAsset.mockResolvedValue({
                status: 200,
                body: mockLogs,
            });

            const { result } = renderHook(() => useLogs("asset-1"), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));

            expect(result.current.data).toEqual(mockLogs);
            expect(mockApiClient.Logs.listLogsByAsset).toHaveBeenCalledWith({
                params: { id: "asset-1" },
                query: {},
            });
        });

        it("should fetch logs with query parameters", async () => {
            const params = { limit: 10, offset: 0 };
            mockApiClient.Logs.listLogsByAsset.mockResolvedValue({
                status: 200,
                body: { data: [], total: 0, limit: 10, offset: 0 },
            });

            const { result } = renderHook(() => useLogs("asset-1", params), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));

            expect(mockApiClient.Logs.listLogsByAsset).toHaveBeenCalledWith({
                params: { id: "asset-1" },
                query: params,
            });
        });

        it("should be disabled when no assetId provided", () => {
            const { result } = renderHook(() => useLogs(""), {
                wrapper: createWrapper(),
            });

            expect(result.current.fetchStatus).toBe("idle");
            expect(mockApiClient.Logs.listLogsByAsset).not.toHaveBeenCalled();
        });
    });

    describe("useLog", () => {
        it("should fetch single log successfully", async () => {
            const mockLog = {
                id: "log-1",
                asset_id: "asset-1",
                user_id: "user1",
                content: "Log content",
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-01T00:00:00Z",
            };

            mockApiClient.Logs.getLogById.mockResolvedValue({
                status: 200,
                body: { data: mockLog },
            });

            const { result } = renderHook(() => useLog("log-1"), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));

            expect(result.current.data).toEqual(mockLog);
            expect(mockApiClient.Logs.getLogById).toHaveBeenCalledWith({
                params: { id: "log-1" },
            });
        });

        it("should be disabled when no ID provided", () => {
            const { result } = renderHook(() => useLog(""), {
                wrapper: createWrapper(),
            });

            expect(result.current.fetchStatus).toBe("idle");
            expect(mockApiClient.Logs.getLogById).not.toHaveBeenCalled();
        });
    });

    describe("useCreateLog", () => {
        it("should create log successfully", async () => {
            const newLog = {
                id: "new-log",
                asset_id: "asset-1",
                user_id: "user1",
                content: "New log",
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-01T00:00:00Z",
            };

            mockApiClient.Logs.createLog.mockResolvedValue({
                status: 201,
                body: { data: newLog },
            });

            const { result } = renderHook(() => useCreateLog(), {
                wrapper: createWrapper(),
            });

            result.current.mutate({ assetId: "asset-1", data: { content: "New log" } });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));

            expect(result.current.data).toEqual({ log: newLog, assetId: "asset-1" });
            expect(mockApiClient.Logs.createLog).toHaveBeenCalledWith({
                params: { id: "asset-1" },
                body: { content: "New log" },
            });
        });
    });

    describe("useUpdateLog", () => {
        it("should update log successfully", async () => {
            const updatedLog = {
                id: "log-1",
                asset_id: "asset-1",
                user_id: "user1",
                content: "Updated log",
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-02T00:00:00Z",
            };

            mockApiClient.Logs.updateLog.mockResolvedValue({
                status: 200,
                body: { data: updatedLog },
            });

            const { result } = renderHook(() => useUpdateLog(), {
                wrapper: createWrapper(),
            });

            result.current.mutate({ id: "log-1", data: { content: "Updated log" } });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));

            expect(result.current.data).toEqual(updatedLog);
            expect(mockApiClient.Logs.updateLog).toHaveBeenCalledWith({
                params: { id: "log-1" },
                body: { content: "Updated log" },
            });
        });
    });

    describe("useDeleteLog", () => {
        it("should delete log successfully", async () => {
            mockApiClient.Logs.deleteLog.mockResolvedValue({
                status: 204,
                body: undefined,
            });

            const { result } = renderHook(() => useDeleteLog(), {
                wrapper: createWrapper(),
            });

            result.current.mutate({ id: "log-1", assetId: "asset-1" });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));

            expect(result.current.data).toEqual({ id: "log-1", assetId: "asset-1" });
            expect(mockApiClient.Logs.deleteLog).toHaveBeenCalledWith({
                params: { id: "log-1" },
            });
        });
    });
});
