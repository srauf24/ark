import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { TagBadge } from "./TagBadge";

describe("TagBadge", () => {
    it("renders tag text", () => {
        render(<TagBadge tag="test-tag" />);
        expect(screen.getByText("test-tag")).toBeInTheDocument();
    });

    it("renders remove button when onRemove is provided", () => {
        render(<TagBadge tag="test-tag" onRemove={() => { }} />);
        expect(screen.getByLabelText("Remove test-tag tag")).toBeInTheDocument();
    });

    it("calls onRemove when remove button is clicked", () => {
        const onRemove = vi.fn();
        render(<TagBadge tag="test-tag" onRemove={onRemove} />);

        fireEvent.click(screen.getByLabelText("Remove test-tag tag"));
        expect(onRemove).toHaveBeenCalled();
    });

    it("does not render remove button when onRemove is missing", () => {
        render(<TagBadge tag="test-tag" />);
        expect(screen.queryByRole("button")).not.toBeInTheDocument();
    });
});
