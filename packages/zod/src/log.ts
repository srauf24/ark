import { z } from "zod";
import { ZBase, ZUuid } from "./common.js";

/**
 * Asset Log Zod schemas matching Go models
 */

// Core AssetLog schema - matches Go model.AssetLog struct
export const ZAssetLog = ZBase.extend({
    asset_id: ZUuid,
    user_id: z.string(),
    content: z.string().min(2).max(10000),
    tags: z.array(z.string().max(50)).max(20).nullable().optional(),
});

// Create Log request - matches Go model.CreateLogRequest
export const ZCreateLogRequest = z.object({
    content: z.string().min(2).max(10000),
    tags: z.array(z.string().max(50)).max(20).optional(),
});

// Update Log request - matches Go model.UpdateLogRequest (all fields optional for PATCH)
export const ZUpdateLogRequest = z.object({
    content: z.string().min(2).max(10000).optional(),
    tags: z.array(z.string().max(50)).max(20).nullable().optional(),
});

// Log query parameters - matches Go model.LogQueryParams
export const ZLogQueryParams = z.object({
    limit: z.coerce.number().int().min(1).max(200).optional(),
    offset: z.coerce.number().int().min(0).optional(),
    tags: z.array(z.string().max(50)).optional(),
    search: z.string().max(100).optional(),
    start_date: z.string().datetime().optional(),
    end_date: z.string().datetime().optional(),
    sort_by: z.enum(["created_at", "updated_at"]).optional(),
    sort_order: z.enum(["asc", "desc"]).optional(),
});

// Log list response - matches Go model.LogListResponse
export const ZLogListResponse = z.object({
    logs: z.array(ZAssetLog),
    total: z.number().int(),
    limit: z.number().int(),
    offset: z.number().int(),
});
