import { z } from "zod";

/**
 * Common Zod schemas for Garden Journal API
 * These schemas are used across multiple resources
 */

// UUID schema - matches Go uuid.UUID
export const ZUuid = z.string().uuid();

// ISO 8601 datetime string - matches Go time.Time serialized as JSON
export const ZTimestamp = z.string().datetime();

// Base model fields - matches Go model.Base
export const ZBase = z.object({
  id: ZUuid,
  created_at: ZTimestamp,
  updated_at: ZTimestamp,
});

// Pagination query parameters
export const ZPaginationQuery = z.object({
  page: z.coerce.number().int().min(1).default(1),
  limit: z.coerce.number().int().min(1).max(100).default(20),
});

// Sorting query parameters
export const ZSortQuery = z.object({
  sort: z.string().optional(),
  order: z.enum(["asc", "desc"]).default("desc"),
});

// Combined query parameters for list endpoints
export const ZListQuery = ZPaginationQuery.merge(ZSortQuery);

// ID path parameter
export const ZIdParam = z.object({
  id: ZUuid,
});

// Common error response
export const ZErrorResponse = z.object({
  message: z.string(),
  code: z.string().optional(),
});

// Success response without data
export const ZSuccessResponse = z.object({
  message: z.string(),
});