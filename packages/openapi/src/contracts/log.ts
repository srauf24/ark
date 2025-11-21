import { getSecurityMetadata } from "../utils.js";
import {
    ZAssetLog,
    ZCreateLogRequest,
    ZUpdateLogRequest,
    ZErrorResponse,
    ZUuid,
} from "@ark/zod";
import { schemaWithPagination } from "@ark/zod";
import { initContract } from "@ts-rest/core";
import { z } from "zod";

const c = initContract();

const metadata = getSecurityMetadata();

export const logContract = c.router(
    {
        listLogsByAsset: {
            summary: "List logs for asset",
            path: "/assets/:id/logs",
            method: "GET",
            description: "Get a paginated list of logs for a specific asset",
            pathParams: z.object({
                id: ZUuid,
            }),
            query: z.object({
                limit: z.coerce.number().int().min(1).max(200).optional(),
                offset: z.coerce.number().int().min(0).optional(),
                tags: z.array(z.string().max(50)).optional(),
                search: z.string().max(100).optional(),
                start_date: z.string().datetime().optional(),
                end_date: z.string().datetime().optional(),
                sort_by: z.enum(["created_at", "updated_at"]).optional(),
                sort_order: z.enum(["asc", "desc"]).optional(),
            }),
            responses: {
                200: schemaWithPagination(ZAssetLog),
                404: ZErrorResponse,
            },
            metadata: metadata,
        },

        createLog: {
            summary: "Create a new log for asset",
            path: "/assets/:id/logs",
            method: "POST",
            description: "Create a new log entry for a specific asset",
            pathParams: z.object({
                id: ZUuid,
            }),
            body: ZCreateLogRequest,
            responses: {
                201: z.object({
                    data: ZAssetLog,
                }),
                400: ZErrorResponse,
                404: ZErrorResponse,
            },
            metadata: metadata,
        },

        getLogById: {
            summary: "Get log by ID",
            path: "/logs/:id",
            method: "GET",
            description: "Get a single log entry by its ID",
            pathParams: z.object({
                id: ZUuid,
            }),
            responses: {
                200: z.object({
                    data: ZAssetLog,
                }),
                404: ZErrorResponse,
            },
            metadata: metadata,
        },

        updateLog: {
            summary: "Update log",
            path: "/logs/:id",
            method: "PATCH",
            description: "Update an existing log entry (partial update)",
            pathParams: z.object({
                id: ZUuid,
            }),
            body: ZUpdateLogRequest,
            responses: {
                200: z.object({
                    data: ZAssetLog,
                }),
                400: ZErrorResponse,
                404: ZErrorResponse,
            },
            metadata: metadata,
        },

        deleteLog: {
            summary: "Delete log",
            path: "/logs/:id",
            method: "DELETE",
            description: "Delete a log entry",
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
