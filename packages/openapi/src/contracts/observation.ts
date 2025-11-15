import { getSecurityMetadata } from "../utils.js";
import {
  ZCreateObservationPayload,
  ZUpdateObservationPayload,
  ZObservation,
  ZErrorResponse,
  ZUuid,
} from "@gardenjournal/zod";
import { schemaWithPagination } from "@gardenjournal/zod";
import { initContract } from "@ts-rest/core";
import z from "zod";

const c = initContract();

const metadata = getSecurityMetadata();

export const observationContract = c.router(
  {
    getObservations: {
      summary: "Get all observations",
      path: "/observations",
      method: "GET",
      description: "Get all observations for the authenticated user",
      query: z.object({
        page: z.coerce.number().int().min(1).optional(),
        limit: z.coerce.number().int().min(1).max(100).optional(),
        sort: z
          .enum(["created_at", "updated_at", "date", "height_cm", "sort_order"])
          .optional(),
        order: z.enum(["asc", "desc"]).optional(),
        search: z.string().min(1).optional(),
      }),
      responses: {
        200: schemaWithPagination(ZObservation),
      },
      metadata: metadata,
    },

    createObservation: {
      summary: "Create a new observation",
      path: "/observations",
      method: "POST",
      description: "Create a new observation for a plant. If date is not provided, it defaults to current timestamp.",
      body: ZCreateObservationPayload,
      responses: {
        201: z.object({
          data: ZObservation,
        }),
        400: ZErrorResponse,
        404: ZErrorResponse,
      },
      metadata: metadata,
    },

    getObservationById: {
      summary: "Get observation by ID",
      path: "/observations/:id",
      method: "GET",
      description: "Get a single observation by ID",
      pathParams: z.object({
        id: ZUuid,
      }),
      responses: {
        200: z.object({
          data: ZObservation,
        }),
        404: ZErrorResponse,
      },
      metadata: metadata,
    },

    updateObservation: {
      summary: "Update observation",
      path: "/observations/:id",
      method: "PUT",
      description: "Update an existing observation",
      pathParams: z.object({
        id: ZUuid,
      }),
      body: ZUpdateObservationPayload,
      responses: {
        200: z.object({
          data: ZObservation,
        }),
        400: ZErrorResponse,
        404: ZErrorResponse,
      },
      metadata: metadata,
    },

    deleteObservation: {
      summary: "Delete observation",
      path: "/observations/:id",
      method: "DELETE",
      description: "Delete an observation",
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