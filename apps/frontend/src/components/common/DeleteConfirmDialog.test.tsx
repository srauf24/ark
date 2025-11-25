import { describe, it, expect, vi } from "vitest";
import { DeleteConfirmDialog } from "./DeleteConfirmDialog";

describe("DeleteConfirmDialog", () => {
    it("should export DeleteConfirmDialog component", () => {
        expect(DeleteConfirmDialog).toBeDefined();
    });

    it("should be a function", () => {
        expect(typeof DeleteConfirmDialog).toBe("function");
    });

    it("should accept required props", () => {
        const props = {
            isOpen: true,
            onConfirm: vi.fn(),
            onCancel: vi.fn(),
            title: "Delete Asset",
            message: "Are you sure?",
        };

        // Verify component can be called with props
        expect(() => DeleteConfirmDialog(props)).not.toThrow();
    });

    it("should accept optional isLoading prop", () => {
        const props = {
            isOpen: true,
            onConfirm: vi.fn(),
            onCancel: vi.fn(),
            title: "Delete Asset",
            message: "Are you sure?",
            isLoading: true,
        };

        expect(() => DeleteConfirmDialog(props)).not.toThrow();
    });
});
