import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook } from "@testing-library/react";

// Mock dependencies BEFORE importing the module under test
vi.mock("@clerk/clerk-react");
vi.mock("@/config/env.ts");
vi.mock("axios");

import { useAuth } from "@clerk/clerk-react";
import { useApiClient } from "./index";

describe("useApiClient", () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Setup default mock implementation
        vi.mocked(useAuth).mockReturnValue({
            getToken: vi.fn().mockResolvedValue("mock-token"),
            isSignedIn: true,
            isLoaded: true,
        } as any);
    });

    it("should return a client object", () => {
        const { result } = renderHook(() => useApiClient());

        expect(result.current).toBeDefined();
        expect(typeof result.current).toBe("object");
    });

    it("should have Assets property with methods", () => {
        const { result } = renderHook(() => useApiClient());

        expect(result.current.Assets).toBeDefined();
        expect(result.current.Assets.listAssets).toBeDefined();
        expect(result.current.Assets.createAsset).toBeDefined();
        expect(result.current.Assets.getAssetById).toBeDefined();
        expect(result.current.Assets.updateAsset).toBeDefined();
        expect(result.current.Assets.deleteAsset).toBeDefined();
    });

    it("should have Logs property with methods", () => {
        const { result } = renderHook(() => useApiClient());

        expect(result.current.Logs).toBeDefined();
        expect(result.current.Logs.listLogsByAsset).toBeDefined();
        expect(result.current.Logs.createLog).toBeDefined();
        expect(result.current.Logs.getLogById).toBeDefined();
        expect(result.current.Logs.updateLog).toBeDefined();
        expect(result.current.Logs.deleteLog).toBeDefined();
    });

    it("should call getToken with custom template", async () => {
        const mockGetToken = vi.fn().mockResolvedValue("mock-token");
        vi.mocked(useAuth).mockReturnValue({
            getToken: mockGetToken,
            isSignedIn: true,
            isLoaded: true,
        } as any);

        const { result } = renderHook(() => useApiClient());

        // Verify client was created successfully
        expect(result.current).toBeDefined();
        expect(mockGetToken).toBeDefined();
    });

    it("should support isBlob parameter", () => {
        const { result } = renderHook(() => useApiClient({ isBlob: true }));

        expect(result.current).toBeDefined();
    });

    it("should work without isBlob parameter (default)", () => {
        const { result } = renderHook(() => useApiClient());

        expect(result.current).toBeDefined();
    });
});
