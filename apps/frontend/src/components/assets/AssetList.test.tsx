// Set environment variables BEFORE any imports
process.env.VITE_CLERK_PUBLISHABLE_KEY = "pk_test_mock_key_for_testing";
process.env.VITE_API_URL = "http://localhost:8080";
process.env.VITE_ENV = "local";

// @vitest-environment happy-dom
import { render, screen, waitFor } from "@testing-library/react";
import { AssetList } from "./AssetList";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter } from "react-router-dom";
import userEvent from "@testing-library/user-event";

// Mock the environment config
vi.mock("@/config/env", () => ({
    API_URL: "http://localhost:3000",
    CLERK_PUBLISHABLE_KEY: "pk_test_mock",
    ENV: "test",
}));

// Mock AssetForm component
vi.mock("./AssetForm", () => ({
    AssetForm: ({ onSuccess, onCancel }: any) => (
        <div data-testid="asset-form">
            <button onClick={onSuccess}>Submit</button>
            <button onClick={onCancel}>Cancel</button>
        </div>
    ),
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
        const spinner = document.querySelector(".animate-spin");
        expect(spinner).toBeInTheDocument();
    });

    it("renders error state on API failure", async () => {
        mockListAssets.mockResolvedValue({ status: 500 });
        renderWithProviders(<AssetList />);

        await waitFor(() => {
            expect(
                screen.getByText("Failed to load assets. Please try again later.")
            ).toBeInTheDocument();
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

    it("opens create dialog when Add Asset button is clicked in empty state", async () => {
        const user = userEvent.setup();
        mockListAssets.mockResolvedValue({
            status: 200,
            body: { assets: [], total: 0 },
        });
        renderWithProviders(<AssetList />);

        await waitFor(() => {
            expect(screen.getByText("Add Asset")).toBeInTheDocument();
        });

        const addButton = screen.getByText("Add Asset");
        await user.click(addButton);

        await waitFor(() => {
            expect(screen.getByTestId("asset-form")).toBeInTheDocument();
            expect(screen.getByText("Create New Asset")).toBeInTheDocument();
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
            expect(screen.getByText("Assets")).toBeInTheDocument();
        });
    });

    it("opens create dialog when Add Asset button is clicked with existing assets", async () => {
        const user = userEvent.setup();
        const mockAssets = [
            {
                id: "1",
                name: "Server 1",
                type: "server",
                created_at: "2023-01-01T00:00:00Z",
                updated_at: "2023-01-01T00:00:00Z",
            },
        ];

        mockListAssets.mockResolvedValue({
            status: 200,
            body: { assets: mockAssets, total: 1 },
        });

        renderWithProviders(<AssetList />);

        await waitFor(() => {
            expect(screen.getByText("Server 1")).toBeInTheDocument();
        });

        const addButton = screen.getByText("Add Asset");
        await user.click(addButton);

        await waitFor(() => {
            expect(screen.getByTestId("asset-form")).toBeInTheDocument();
            expect(screen.getByText("Create New Asset")).toBeInTheDocument();
        });
    });
});
