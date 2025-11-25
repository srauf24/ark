import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { LogList } from "./LogList";
import { createWrapper } from "../../../test/utils";

// Mock hooks
const mockUseLogs = vi.fn();
const mockDeleteLog = vi.fn();

vi.mock("@/hooks/useLogs", () => ({
    useLogs: () => mockUseLogs(),
    useDeleteLog: () => ({
        mutateAsync: mockDeleteLog,
        isPending: false,
    }),
    useCreateLog: () => ({ mutateAsync: vi.fn(), isPending: false }),
    useUpdateLog: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

describe("LogList", () => {
    const assetId = "asset-1";

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders loading state", () => {
        mockUseLogs.mockReturnValue({ isLoading: true });
        render(<LogList assetId={assetId} />, { wrapper: createWrapper() });
        expect(screen.getAllByRole("status")).toHaveLength(1); // Skeleton has role="status" implicitly or we check for class
        // Or check for absence of "Add Log" button or presence of skeletons
        // Since Skeleton doesn't have a specific role by default in shadcn, we can check for structure
        // But let's just check that "Logs" header is present
        expect(screen.getByText("Logs")).toBeInTheDocument();
    });

    it("renders error state", () => {
        mockUseLogs.mockReturnValue({ isError: true, refetch: vi.fn() });
        render(<LogList assetId={assetId} />, { wrapper: createWrapper() });
        expect(screen.getByText("Failed to load logs.")).toBeInTheDocument();
    });

    it("renders empty state", () => {
        mockUseLogs.mockReturnValue({ data: { logs: [] }, isLoading: false });
        render(<LogList assetId={assetId} />, { wrapper: createWrapper() });
        expect(screen.getByText("No logs recorded yet.")).toBeInTheDocument();
        expect(screen.getByRole("button", { name: "Create your first log" })).toBeInTheDocument();
    });

    it("renders list of logs", () => {
        const mockLogs = [
            {
                id: "log-1",
                asset_id: assetId,
                content: "Log 1",
                created_at: new Date().toISOString(),
            },
            {
                id: "log-2",
                asset_id: assetId,
                content: "Log 2",
                created_at: new Date().toISOString(),
            },
        ];
        mockUseLogs.mockReturnValue({ data: { logs: mockLogs }, isLoading: false });

        render(<LogList assetId={assetId} />, { wrapper: createWrapper() });

        expect(screen.getByText("Log 1")).toBeInTheDocument();
        expect(screen.getByText("Log 2")).toBeInTheDocument();
    });

    it("opens create modal", () => {
        mockUseLogs.mockReturnValue({ data: { logs: [] }, isLoading: false });
        render(<LogList assetId={assetId} />, { wrapper: createWrapper() });

        fireEvent.click(screen.getByRole("button", { name: "Add Log" }));
        expect(screen.getByText("Add Log")).toBeInTheDocument(); // Modal title
    });
});
