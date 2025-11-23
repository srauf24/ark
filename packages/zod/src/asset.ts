import { z } from "zod";
import { ZBase, ZUuid } from "./common.js";

/**
 * Asset Zod schemas matching Go models
 */

// Asset type enum
export const ZAssetType = z.enum(["server", "vm", "nas", "container", "network", "other"]);

// Asset metadata - flexible JSON structure for asset-specific details
// Examples: CPU specs, RAM, IP addresses, ports, etc.
export const ZAssetMetadata = z.record(z.any());

// Core Asset schema - matches Go model.Asset struct
export const ZAsset = ZBase.extend({
    user_id: z.string(),
    name: z.string().min(1).max(100),
    type: ZAssetType.nullable().optional(),
    hostname: z.string().max(255).nullable().optional(),
    metadata: ZAssetMetadata.nullable().optional(),
});

export type Asset = z.infer<typeof ZAsset>;

// Create Asset request - matches Go model.CreateAssetRequest
export const ZCreateAssetRequest = z.object({
    name: z.string().min(1).max(100),
    type: ZAssetType.optional(),
    hostname: z.string().max(255).optional(),
    metadata: ZAssetMetadata.optional(),
});

// Update Asset request - matches Go model.UpdateAssetRequest (all fields optional for PATCH)
export const ZUpdateAssetRequest = z.object({
    name: z.string().min(1).max(100).optional(),
    type: ZAssetType.optional(),
    hostname: z.string().max(255).optional(),
    metadata: ZAssetMetadata.optional(),
});

// Asset query parameters - matches Go model.AssetQueryParams
export const ZAssetQueryParams = z.object({
    limit: z.coerce.number().int().min(1).max(100).optional(),
    offset: z.coerce.number().int().min(0).optional(),
    type: ZAssetType.optional(),
    search: z.string().max(100).optional(),
    sort_by: z.enum(["name", "created_at", "updated_at"]).optional(),
    sort_order: z.enum(["asc", "desc"]).optional(),
});

// Asset list response - matches Go model.AssetListResponse
export const ZAssetListResponse = z.object({
    assets: z.array(ZAsset),
    total: z.number().int(),
    limit: z.number().int(),
    offset: z.number().int(),
});
