import { describe, it, expect, vi } from "vitest";

// Mock dependencies BEFORE importing the module under test
const mockGetToken = vi.fn().mockResolvedValue("mock-token");

vi.mock("@clerk/clerk-react", () => ({
    useAuth: () => ({
        getToken: mockGetToken,
        isSignedIn: true,
        isLoaded: true,
    }),
}));

vi.mock("@/config/env.ts", () => ({
    API_URL: "http://localhost:8080",
}));

vi.mock("axios", () => ({
    default: {
        request: vi.fn(),
    },
    isAxiosError: vi.fn(),
}));

import { useApiClient } from "./index";

describe("useApiClient", () => {
    it("should return a client object", () => {
        const client = useApiClient();

        expect(client).toBeDefined();
        expect(typeof client).toBe("object");
    });

    it("should have Assets property with methods", () => {
        const client = useApiClient();

        expect(client.Assets).toBeDefined();
        expect(client.Assets.listAssets).toBeDefined();
        expect(client.Assets.createAsset).toBeDefined();
        expect(client.Assets.getAssetById).toBeDefined();
        expect(client.Assets.updateAsset).toBeDefined();
        expect(client.Assets.deleteAsset).toBeDefined();
    });

    it("should have Logs property with methods", () => {
        const client = useApiClient();

        expect(client.Logs).toBeDefined();
        expect(client.Logs.listLogsByAsset).toBeDefined();
        expect(client.Logs.createLog).toBeDefined();
        expect(client.Logs.getLogById).toBeDefined();
        expect(client.Logs.updateLog).toBeDefined();
        expect(client.Logs.deleteLog).toBeDefined();
    });

    it("should support isBlob parameter", () => {
        const client = useApiClient({ isBlob: true });

        expect(client).toBeDefined();
        expect(client.Assets).toBeDefined();
    });

    it("should work without isBlob parameter (default)", () => {
        const client = useApiClient();

        expect(client).toBeDefined();
        expect(client.Assets).toBeDefined();
    });
});
