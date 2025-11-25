import { renderHook, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { useLogs, useCreateLog, useUpdateLog, useDeleteLog } from "./useLogs";
import { createWrapper } from "../../test/utils";

// Mock API client
const mockListLogsByAsset = vi.fn();
const mockCreateLog = vi.fn();
const mockUpdateLog = vi.fn();
const mockDeleteLog = vi.fn();

vi.mock("@/api", () => ({
    useApiClient: () => ({
        logs: {
            listLogsByAsset: mockListLogsByAsset,
            createLog: mockCreateLog,
            updateLog: mockUpdateLog,
            deleteLog: mockDeleteLog,
        },
    }),
}));

describe("useLogs Hooks", () => {
    const assetId = "asset-123";

    beforeEach(() => {
        vi.clearAllMocks();
    });

    describe("useLogs", () => {
        it("fetches logs successfully", async () => {
            const mockData = { logs: [], total: 0, limit: 50, offset: 0 };
            mockListLogsByAsset.mockResolvedValue({ status: 200, body: mockData });

            const { result } = renderHook(() => useLogs(assetId), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));
            expect(result.current.data).toEqual(mockData);
            expect(mockListLogsByAsset).toHaveBeenCalledWith({
                params: { id: assetId },
                query: {},
            });
        });
    });

    describe("useCreateLog", () => {
        it("creates log successfully", async () => {
            const mockLog = { id: "log-1", content: "test" };
            mockCreateLog.mockResolvedValue({ status: 201, body: mockLog });

            const { result } = renderHook(() => useCreateLog(assetId), {
                wrapper: createWrapper(),
            });

            result.current.mutate({ content: "test" });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));
            expect(mockCreateLog).toHaveBeenCalledWith({
                params: { id: assetId },
                body: { content: "test" },
            });
        });
    });

    describe("useUpdateLog", () => {
        it("updates log successfully", async () => {
            const mockLog = { id: "log-1", content: "updated" };
            mockUpdateLog.mockResolvedValue({ status: 200, body: mockLog });

            const { result } = renderHook(() => useUpdateLog(assetId), {
                wrapper: createWrapper(),
            });

            result.current.mutate({ id: "log-1", data: { content: "updated" } });

            await waitFor(() => expect(result.current.isSuccess).toBe(true));
            expect(mockUpdateLog).toHaveBeenCalledWith({
                params: { id: "log-1" },
                body: { content: "updated" },
            });
        });
    });

    describe("useDeleteLog", () => {
        it("deletes log successfully", async () => {
            mockDeleteLog.mockResolvedValue({ status: 204, body: null });

            const { result } = renderHook(() => useDeleteLog(assetId), {
                wrapper: createWrapper(),
            });

            result.current.mutate("log-1");

            await waitFor(() => expect(result.current.isSuccess).toBe(true));
            expect(mockDeleteLog).toHaveBeenCalledWith({
                params: { id: "log-1" },
            });
        });
    });
});
