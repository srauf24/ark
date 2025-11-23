// @vitest-environment happy-dom
import { render, screen, waitFor } from "@testing-library/react";
import { AssetList } from "./AssetList";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter } from "react-router-dom";

// Mock the environment config
vi.mock("@/config/env", () => ({
    API_URL: "http://localhost:3000",
    CLERK_PUBLISHABLE_KEY: "pk_test_mock",
    ENV: "test",
}));

// Mock the API client hook
const mockListAssets = vi.fn();

vi.mock("@/api", () => ({
    useApiClient: () => ({
        Assets: {
            listAssets: mockListAssets,
        },
    }),
}));

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: false,
        },
    },
});

const renderWithProviders = (component: React.ReactNode) => {
    return render(
        <QueryClientProvider client={queryClient}>
            <BrowserRouter>{component}</BrowserRouter>
        </QueryClientProvider>
    );
};

describe("AssetList", () => {
    beforeEach(() => {
        vi.clearAllMocks();
        queryClient.clear();
    });

    it("renders loading state initially", () => {
        mockListAssets.mockReturnValue(new Promise(() => { })); // Never resolves
        renderWithProviders(<AssetList />);
        expect(screen.getByRole("status", { hidden: true })).toBeInTheDocument(); // Loader2 implies status role usually, but let's check class or existence
        // Since Loader2 is an SVG, we can check for its presence via class or container
        // A better way is to check for the container or specific accessible element if added. 
        // For now, let's check if the container with the spinner class exists
        const spinner = document.querySelector(".animate-spin");
        expect(spinner).toBeInTheDocument();
    });

    it("renders error state on API failure", async () => {
        mockListAssets.mockResolvedValue({ status: 500 });
        renderWithProviders(<AssetList />);

        await waitFor(() => {
            expect(screen.getByText("Failed to load assets. Please try again later.")).toBeInTheDocument();
        });
    });

    it("renders empty state when no assets", async () => {
        mockListAssets.mockResolvedValue({
            status: 200,
            body: { assets: [], total: 0 },
        });
        renderWithProviders(<AssetList />);

        await waitFor(() => {
            expect(screen.getByText("No assets found")).toBeInTheDocument();
            expect(screen.getByText("Add Asset")).toBeInTheDocument();
        });
    });

    it("renders list of assets when data exists", async () => {
        const mockAssets = [
            {
                id: "1",
                name: "Server 1",
                type: "server",
                created_at: "2023-01-01T00:00:00Z",
                updated_at: "2023-01-01T00:00:00Z",
            },
            {
                id: "2",
                name: "NAS 1",
                type: "nas",
                created_at: "2023-01-01T00:00:00Z",
                updated_at: "2023-01-01T00:00:00Z",
            },
        ];

        mockListAssets.mockResolvedValue({
            status: 200,
            body: { assets: mockAssets, total: 2 },
        });

        renderWithProviders(<AssetList />);

        await waitFor(() => {
            expect(screen.getByText("Server 1")).toBeInTheDocument();
            expect(screen.getByText("NAS 1")).toBeInTheDocument();
        });
    });
});
