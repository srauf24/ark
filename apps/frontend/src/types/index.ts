import { z } from "zod";
import {
    ZAsset,
    ZAssetLog,
    ZCreateAssetRequest,
    ZUpdateAssetRequest,
    ZCreateLogRequest,
    ZUpdateLogRequest,
    ZAssetQueryParams,
    ZLogQueryParams,
} from "@ark/zod";

/**
 * Domain Models
 */
export type Asset = z.infer<typeof ZAsset>;
export type AssetLog = z.infer<typeof ZAssetLog>;

/**
 * Request DTOs
 */
export type CreateAssetRequest = z.infer<typeof ZCreateAssetRequest>;
export type UpdateAssetRequest = z.infer<typeof ZUpdateAssetRequest>;
export type CreateLogRequest = z.infer<typeof ZCreateLogRequest>;
export type UpdateLogRequest = z.infer<typeof ZUpdateLogRequest>;

/**
 * Query Parameters
 */
export type AssetQueryParams = z.infer<typeof ZAssetQueryParams>;
export type LogQueryParams = z.infer<typeof ZLogQueryParams>;

/**
 * Pagination Types
 */
export interface PaginationParams {
    limit?: number;
    offset?: number;
}

export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    limit: number;
    offset: number;
}

/**
 * API Error Response
 */
export interface ApiError {
    error: string;
    details?: any;
}

/**
 * Asset Types Enum
 */
export type AssetType = "server" | "vm" | "nas" | "container" | "network" | "other";
