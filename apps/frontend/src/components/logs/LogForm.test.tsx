import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { LogForm } from "./LogForm";

describe("LogForm", () => {
    it("renders form fields", () => {
        render(<LogForm onSubmit={() => { }} isSubmitting={false} />);
        expect(screen.getByLabelText("Content")).toBeInTheDocument();
        expect(screen.getByLabelText("Tags")).toBeInTheDocument();
        expect(screen.getByRole("button", { name: "Save Log" })).toBeInTheDocument();
    });

    it("validates required fields", async () => {
        render(<LogForm onSubmit={() => { }} isSubmitting={false} />);

        fireEvent.click(screen.getByRole("button", { name: "Save Log" }));

        await waitFor(() => {
            expect(screen.getByText(/String must contain at least 2 character/i)).toBeInTheDocument();
        });
    });

    it("adds tags on Enter", async () => {
        render(<LogForm onSubmit={() => { }} isSubmitting={false} />);

        const tagInput = screen.getByLabelText("Tags");
        fireEvent.change(tagInput, { target: { value: "new-tag" } });
        fireEvent.keyDown(tagInput, { key: "Enter" });

        expect(screen.getByText("new-tag")).toBeInTheDocument();
        expect(tagInput).toHaveValue("");
    });

    it("submits valid data", async () => {
        const onSubmit = vi.fn();
        render(<LogForm onSubmit={onSubmit} isSubmitting={false} />);

        fireEvent.change(screen.getByLabelText("Content"), { target: { value: "Test content" } });

        const tagInput = screen.getByLabelText("Tags");
        fireEvent.change(tagInput, { target: { value: "tag1" } });
        fireEvent.keyDown(tagInput, { key: "Enter" });

        fireEvent.click(screen.getByRole("button", { name: "Save Log" }));

        await waitFor(() => {
            expect(onSubmit).toHaveBeenCalledWith({
                content: "Test content",
                tags: ["tag1"],
            });
        });
    });
});
