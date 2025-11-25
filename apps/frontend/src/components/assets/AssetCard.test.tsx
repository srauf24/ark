// @vitest-environment happy-dom
import { render, screen } from "@testing-library/react";
import { AssetCard } from "./AssetCard";
import { BrowserRouter } from "react-router-dom";
import { describe, it, expect } from "vitest";
import { type Asset } from "@ark/zod";

const mockAsset: Asset = {
    id: "123e4567-e89b-12d3-a456-426614174000",
    user_id: "user_123",
    name: "Test Server",
    type: "server",
    hostname: "server.local",
    metadata: {},
    created_at: "2023-01-01T00:00:00Z",
    updated_at: "2023-01-02T00:00:00Z",
};

describe("AssetCard", () => {
    it("renders asset name and hostname", () => {
        render(
            <BrowserRouter>
                <AssetCard asset={mockAsset} />
            </BrowserRouter>
        );

        expect(screen.getByText("Test Server")).toBeInTheDocument();
        expect(screen.getByText("server.local")).toBeInTheDocument();
    });

    it("renders correct type badge", () => {
        render(
            <BrowserRouter>
                <AssetCard asset={mockAsset} />
            </BrowserRouter>
        );

        expect(screen.getByText("server")).toBeInTheDocument();
    });

    it("renders updated date correctly", () => {
        render(
            <BrowserRouter>
                <AssetCard asset={mockAsset} />
            </BrowserRouter>
        );

        expect(screen.getByText("Jan 2, 2023")).toBeInTheDocument();
    });

    it("links to correct asset detail page", () => {
        render(
            <BrowserRouter>
                <AssetCard asset={mockAsset} />
            </BrowserRouter>
        );

        const link = screen.getByRole("link");
        expect(link).toHaveAttribute("href", `/assets/${mockAsset.id}`);
    });

    it("renders default icon for unknown type", () => {
        const unknownAsset = { ...mockAsset, type: "unknown" as any };
        render(
            <BrowserRouter>
                <AssetCard asset={unknownAsset} />
            </BrowserRouter>
        );

        // HelpCircle icon is used for default/unknown
        // We can check if the badge renders "unknown" (or whatever the type is)
        expect(screen.getByText("unknown")).toBeInTheDocument();
    });
});
