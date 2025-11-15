import { getSecurityMetadata } from "../utils.js";
import {
  ZCreatePlantPayload,
  ZUpdatePlantPayload,
  ZPlant,
  ZPopulatedPlant,
  ZErrorResponse,
  ZUuid,
  ZTimestamp,
} from "@gardenjournal/zod";
import { schemaWithPagination } from "@gardenjournal/zod";
import { initContract } from "@ts-rest/core";
import z from "zod";

const c = initContract();

const metadata = getSecurityMetadata();

export const plantContract = c.router(
  {
    getPlants: {
      summary: "Get all plants",
      path: "/plants",
      method: "GET",
      description: "Get all plants for the authenticated user",
      query: z.object({
        page: z.coerce.number().int().min(1).optional(),
        limit: z.coerce.number().int().min(1).max(100).optional(),
        sort: z
          .enum([
            "created_at",
            "updated_at",
            "name",
            "species",
            "location",
            "planted_date",
            "sort_order",
          ])
          .optional(),
        order: z.enum(["asc", "desc"]).optional(),
        search: z.string().min(1).optional(),
        species: z.string().min(1).optional(),
        location: z.string().min(1).optional(),
        plantedFrom: ZTimestamp.optional(),
        plantedTo: ZTimestamp.optional(),
      }),
      responses: {
        200: schemaWithPagination(ZPopulatedPlant),
      },
      metadata: metadata,
    },

    createPlant: {
      summary: "Create a new plant",
      path: "/plants",
      method: "POST",
      description: "Create a new plant for the authenticated user",
      body: ZCreatePlantPayload,
      responses: {
        201: z.object({
          data: ZPlant,
        }),
        400: ZErrorResponse,
      },
      metadata: metadata,
    },

    getPlantById: {
      summary: "Get plant by ID",
      path: "/plants/:id",
      method: "GET",
      description: "Get a single plant by ID with populated observations",
      pathParams: z.object({
        id: ZUuid,
      }),
      responses: {
        200: z.object({
          data: ZPopulatedPlant,
        }),
        404: ZErrorResponse,
      },
      metadata: metadata,
    },

    updatePlant: {
      summary: "Update plant",
      path: "/plants/:id",
      method: "PUT",
      description: "Update an existing plant",
      pathParams: z.object({
        id: ZUuid,
      }),
      body: ZUpdatePlantPayload,
      responses: {
        200: z.object({
          data: ZPlant,
        }),
        400: ZErrorResponse,
        404: ZErrorResponse,
      },
      metadata: metadata,
    },

    deletePlant: {
      summary: "Delete plant",
      path: "/plants/:id",
      method: "DELETE",
      description: "Delete a plant and all its observations",
      pathParams: z.object({
        id: ZUuid,
      }),
      responses: {
        204: z.void(),
      },
      metadata: metadata,
    },
  },
  {
    pathPrefix: "/api/v1",
  }
);