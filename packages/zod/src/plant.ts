import { z } from "zod";
import { ZBase, ZListQuery, ZTimestamp, ZUuid } from "./common.js";
import { schemaWithPagination } from "./utils.js";

/**
 * Plant Zod schemas matching Go models
 */

// Plant metadata - flexible JSON structure for additional plant data
export const ZPlantMetadata = z.object({
  tags: z.array(z.string()).optional(),
  wateringFrequency: z.string().optional(),
  lastWateredAt: ZTimestamp.optional(),
  sunlightLevel: z.string().optional(),
  soilType: z.string().optional(),
  potSizeCm: z.number().optional(),
  fertilizerType: z.string().optional(),
  lastFertilizedAt: ZTimestamp.optional(),
  lastWeatherSnapshotId: ZUuid.optional(),
  averageTempC: z.number().optional(),
  averageSunshineHrs: z.number().optional(),
  healthStatus: z.string().optional(),
  growthStage: z.string().optional(),
  heightCm: z.number().optional(),
  colorTag: z.string().optional(),
  imageUrl: z.string().url().optional(),
  emojiIcon: z.string().optional(),
  aiInsightSummary: z.string().optional(),
}).passthrough();

// Core Plant schema - matches Go plant.Plant struct
export const ZPlant = ZBase.extend({
  userId: z.string(),
  name: z.string().min(1).max(100),
  species: z.string().min(1).max(100),
  location: z.string().max(255).nullable(),
  plantedDate: ZTimestamp.nullable(),
  notes: z.string().max(1000).nullable(),
  metadata: ZPlantMetadata.nullable(),
  sortOrder: z.number().int(),
});

// Populated Plant - includes observations
export const ZPopulatedPlant = ZPlant.extend({
  observations: z.array(z.any()),
});

// Create Plant payload - matches Go plant.CreatePlantPayload
export const ZCreatePlantPayload = z.object({
  name: z.string().min(1).max(100),
  species: z.string().min(1).max(100),
  location: z.string().max(255).optional(),
  plantedDate: ZTimestamp.optional(),
  notes: z.string().max(1000).optional(),
  metadata: ZPlantMetadata.optional(),
});

// Update Plant payload - matches Go plant.UpdatePlantPayload (all fields optional)
export const ZUpdatePlantPayload = z.object({
  name: z.string().min(1).max(100).optional(),
  species: z.string().min(1).max(100).optional(),
  location: z.string().max(255).optional(),
  plantedDate: ZTimestamp.optional(),
  notes: z.string().max(1000).optional(),
  metadata: ZPlantMetadata.optional(),
});

// Get Plants query - matches Go plant.GetPlantsQuery
export const ZGetPlantsQuery = ZListQuery.extend({
  search: z.string().min(1).optional(),
  species: z.string().min(1).optional(),
  location: z.string().min(1).optional(),
  plantedFrom: ZTimestamp.optional(),
  plantedTo: ZTimestamp.optional(),
  sort: z.enum(["created_at", "updated_at", "name", "species", "location", "planted_date", "sort_order"]).optional(),
});

// Response schemas
export const ZPlantResponse = z.object({
  data: ZPlant,
});

export const ZPopulatedPlantResponse = z.object({
  data: ZPopulatedPlant,
});

export const ZPlantsListResponse = schemaWithPagination(ZPopulatedPlant);
