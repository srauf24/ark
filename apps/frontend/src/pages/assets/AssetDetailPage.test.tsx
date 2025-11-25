// @vitest-environment happy-dom
import { render, screen, waitFor } from "@testing-library/react";
import { AssetDetailPage } from "./AssetDetailPage";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MemoryRouter, Routes, Route } from "react-router-dom";

// Mock the environment config
vi.mock("@/config/env", () => ({
    API_URL: "http://localhost:3000",
    CLERK_PUBLISHABLE_KEY: "pk_test_mock",
    ENV: "test",
}));

// Mock the API client hook
const mockGetAssetById = vi.fn();

vi.mock("@/api", () => ({
    useApiClient: () => ({
        Assets: {
            getAssetById: mockGetAssetById,
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

const renderWithProviders = (component: React.ReactNode, initialEntries = ["/assets/1"]) => {
    return render(
        <QueryClientProvider client={queryClient}>
            <MemoryRouter initialEntries={initialEntries}>
                <Routes>
                    <Route path="/assets/:id" element={component} />
                </Routes>
            </MemoryRouter>
        </QueryClientProvider>
    );
};

describe("AssetDetailPage", () => {
    beforeEach(() => {
        vi.clearAllMocks();
        queryClient.clear();
    });

    it("renders loading state initially", () => {
        mockGetAssetById.mockReturnValue(new Promise(() => { })); // Never resolves
        renderWithProviders(<AssetDetailPage />);
        expect(screen.getByRole("status")).toBeInTheDocument();
    });

    it("renders error state on API failure", async () => {
        mockGetAssetById.mockResolvedValue({ status: 500 });
        renderWithProviders(<AssetDetailPage />);

        await waitFor(() => {
            expect(screen.getByText("Failed to load asset details. The asset may not exist or you don't have permission to view it.")).toBeInTheDocument();
        });
    });

    it("renders asset details when data exists", async () => {
        const mockAsset = {
            id: "1",
            name: "Test Server",
            type: "server",
            hostname: "server.local",
            metadata: { cpu: "4 cores", ram: "16GB" },
            created_at: "2023-01-01T00:00:00Z",
            updated_at: "2023-01-01T00:00:00Z",
        };

        mockGetAssetById.mockResolvedValue({
            status: 200,
            body: { data: mockAsset },
        });

        renderWithProviders(<AssetDetailPage />);

        await waitFor(() => {
            expect(screen.getByText("Test Server")).toBeInTheDocument();
            expect(screen.getByText("server.local")).toBeInTheDocument();
            expect(screen.getByText("server")).toBeInTheDocument();
            expect(screen.getByText(/4 cores/)).toBeInTheDocument();
        });
    });
});
