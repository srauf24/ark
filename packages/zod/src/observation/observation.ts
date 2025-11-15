import { z } from "zod";
import { ZBase, ZListQuery, ZTimestamp, ZUuid } from "../common.js";
import { schemaWithPagination } from "../utils.js";

/**
 * Observation Zod schemas matching Go models
 */

// Core Observation schema - matches Go observation.Observation struct
export const ZObservation = ZBase.extend({
  userId: z.string(),
  plantId: ZUuid,
  date: ZTimestamp,
  heightCm: z.number().nullable(),
  notes: z.string().nullable(),
  sortOrder: z.number().int(),
});

// Create Observation payload - matches Go observation.CreateObservationPayload
export const ZCreateObservationPayload = z.object({
  plantId: ZUuid,
  date: ZTimestamp.optional(),
  heightCm: z.number().optional(),
  notes: z.string().optional(),
});

// Update Observation payload - matches Go observation.UpdateObservationPayload (all fields optional)
export const ZUpdateObservationPayload = z.object({
  heightCm: z.number().optional(),
  notes: z.string().optional(),
});

// Get Observations query - matches Go observation.GetObservationsQuery
export const ZGetObservationsQuery = ZListQuery.extend({
  search: z.string().min(1).optional(),
  sort: z.enum(["created_at", "updated_at", "date", "height_cm", "sort_order"]).optional(),
});

// Response schemas
export const ZObservationResponse = z.object({
  data: ZObservation,
});

export const ZObservationsListResponse = schemaWithPagination(ZObservation);
