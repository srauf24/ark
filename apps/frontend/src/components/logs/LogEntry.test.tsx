import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { LogEntry } from "./LogEntry";

describe("LogEntry", () => {
    const mockLog = {
        id: "log-1",
        asset_id: "asset-1",
        user_id: "user-1",
        content: "Test log content",
        tags: ["tag1", "tag2"],
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
    };

    it("renders log content and tags", () => {
        render(
            <LogEntry
                log={mockLog}
                onEdit={() => { }}
                onDelete={() => { }}
            />
        );

        expect(screen.getByText("Test log content")).toBeInTheDocument();
        expect(screen.getByText("tag1")).toBeInTheDocument();
        expect(screen.getByText("tag2")).toBeInTheDocument();
    });

    it("calls onEdit when edit button is clicked", () => {
        const onEdit = vi.fn();
        render(
            <LogEntry
                log={mockLog}
                onEdit={onEdit}
                onDelete={() => { }}
            />
        );

        fireEvent.click(screen.getByLabelText("Edit log"));
        expect(onEdit).toHaveBeenCalled();
    });

    it("calls onDelete when delete button is clicked", () => {
        const onDelete = vi.fn();
        render(
            <LogEntry
                log={mockLog}
                onEdit={() => { }}
                onDelete={onDelete}
            />
        );

        fireEvent.click(screen.getByLabelText("Delete log"));
        expect(onDelete).toHaveBeenCalled();
    });
});
