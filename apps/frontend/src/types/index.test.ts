import { describe, it, expect } from "vitest";
import type {
    Asset,
    AssetLog,
    CreateAssetRequest,
    UpdateAssetRequest,
    CreateLogRequest,
    UpdateLogRequest,
    AssetQueryParams,
    LogQueryParams,
    PaginationParams,
    PaginatedResponse,
    ApiError,
    AssetType,
} from "./index";

describe("TypeScript Types", () => {
    describe("Asset Types", () => {
        it("should have correct Asset type structure", () => {
            const asset: Asset = {
                id: "123e4567-e89b-12d3-a456-426614174000",
                user_id: "user_123",
                name: "Test Server",
                type: "server",
                hostname: "test.example.com",
                metadata: { cpu: "Intel i7", ram: "16GB" },
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-01T00:00:00Z",
            };

            expect(asset).toBeDefined();
            expect(asset.id).toBe("123e4567-e89b-12d3-a456-426614174000");
            expect(asset.name).toBe("Test Server");
        });

        it("should have correct AssetLog type structure", () => {
            const log: AssetLog = {
                id: "123e4567-e89b-12d3-a456-426614174001",
                asset_id: "123e4567-e89b-12d3-a456-426614174000",
                user_id: "user_123",
                content: "Server maintenance completed",
                tags: ["maintenance", "server"],
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-01T00:00:00Z",
            };

            expect(log).toBeDefined();
            expect(log.content).toBe("Server maintenance completed");
        });
    });

    describe("Request DTO Types", () => {
        it("should have correct CreateAssetRequest type", () => {
            const request: CreateAssetRequest = {
                name: "New Server",
                type: "server",
                hostname: "new.example.com",
                metadata: { location: "datacenter-1" },
            };

            expect(request).toBeDefined();
            expect(request.name).toBe("New Server");
        });

        it("should have correct UpdateAssetRequest type with all optional fields", () => {
            const request: UpdateAssetRequest = {
                name: "Updated Server",
            };

            expect(request).toBeDefined();
            expect(request.name).toBe("Updated Server");
        });

        it("should have correct CreateLogRequest type", () => {
            const request: CreateLogRequest = {
                content: "New log entry",
                tags: ["info"],
            };

            expect(request).toBeDefined();
            expect(request.content).toBe("New log entry");
        });

        it("should have correct UpdateLogRequest type with all optional fields", () => {
            const request: UpdateLogRequest = {
                content: "Updated log entry",
            };

            expect(request).toBeDefined();
            expect(request.content).toBe("Updated log entry");
        });
    });

    describe("Query Parameter Types", () => {
        it("should have correct AssetQueryParams type", () => {
            const params: AssetQueryParams = {
                limit: 50,
                offset: 0,
                type: "server",
                search: "test",
                sort_by: "name",
                sort_order: "asc",
            };

            expect(params).toBeDefined();
            expect(params.limit).toBe(50);
        });

        it("should have correct LogQueryParams type", () => {
            const params: LogQueryParams = {
                limit: 100,
                offset: 0,
                tags: ["maintenance"],
                search: "server",
                start_date: "2024-01-01T00:00:00Z",
                end_date: "2024-12-31T23:59:59Z",
                sort_by: "created_at",
                sort_order: "desc",
            };

            expect(params).toBeDefined();
            expect(params.limit).toBe(100);
        });
    });

    describe("Pagination Types", () => {
        it("should have correct PaginationParams type", () => {
            const params: PaginationParams = {
                limit: 25,
                offset: 50,
            };

            expect(params).toBeDefined();
            expect(params.limit).toBe(25);
        });

        it("should have correct PaginatedResponse type", () => {
            const response: PaginatedResponse<Asset> = {
                data: [],
                total: 100,
                limit: 25,
                offset: 0,
            };

            expect(response).toBeDefined();
            expect(response.total).toBe(100);
        });
    });

    describe("Error Types", () => {
        it("should have correct ApiError type", () => {
            const error: ApiError = {
                error: "Not found",
                details: { resource: "asset", id: "123" },
            };

            expect(error).toBeDefined();
            expect(error.error).toBe("Not found");
        });
    });

    describe("Enum Types", () => {
        it("should have correct AssetType values", () => {
            const types: AssetType[] = ["server", "vm", "nas", "container", "network", "other"];

            types.forEach((type) => {
                const asset: Partial<Asset> = { type };
                expect(asset.type).toBe(type);
            });
        });
    });
});
