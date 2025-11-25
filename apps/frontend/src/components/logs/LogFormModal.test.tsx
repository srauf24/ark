import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { LogFormModal } from "./LogFormModal";
import { createWrapper } from "../../../test/utils";

// Mock hooks
const mockMutateAsync = vi.fn();
vi.mock("@/hooks/useLogs", () => ({
    useCreateLog: () => ({
        mutateAsync: mockMutateAsync,
        isPending: false,
    }),
    useUpdateLog: () => ({
        mutateAsync: mockMutateAsync,
        isPending: false,
    }),
}));

describe("LogFormModal", () => {
    const assetId = "asset-1";

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders Add Log title when no initialData", () => {
        render(
            <LogFormModal
                open={true}
                onOpenChange={() => { }}
                assetId={assetId}
            />,
            { wrapper: createWrapper() }
        );
        expect(screen.getByText("Add Log")).toBeInTheDocument();
    });

    it("renders Edit Log title when initialData provided", () => {
        const mockLog = {
            id: "log-1",
            asset_id: assetId,
            user_id: "user-1",
            content: "Test content",
            tags: ["tag1"],
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };

        render(
            <LogFormModal
                open={true}
                onOpenChange={() => { }}
                initialData={mockLog}
                assetId={assetId}
            />,
            { wrapper: createWrapper() }
        );
        expect(screen.getByText("Edit Log")).toBeInTheDocument();
        expect(screen.getByDisplayValue("Test content")).toBeInTheDocument();
    });

    it("calls create mutation on submit", async () => {
        render(
            <LogFormModal
                open={true}
                onOpenChange={() => { }}
                assetId={assetId}
            />,
            { wrapper: createWrapper() }
        );

        fireEvent.change(screen.getByLabelText("Content"), { target: { value: "New log" } });
        fireEvent.click(screen.getByRole("button", { name: "Save Log" }));

        await waitFor(() => {
            expect(mockMutateAsync).toHaveBeenCalledWith({
                content: "New log",
                tags: [],
            });
        });
    });
});
