import { getSecurityMetadata } from "../utils.js";
import {
    ZAssetLog,
    ZCreateLogRequest,
    ZUpdateLogRequest,
    ZLogListResponse,
    ZLogQueryParams,
    ZErrorResponse,
    ZUuid,
} from "@gardenjournal/zod";
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
            pathParams: z.object({ id: ZUuid }),
            query: ZLogQueryParams,
            responses: {
                200: ZLogListResponse,
                404: ZErrorResponse,
            },
            metadata: metadata,
        },

        createLog: {
            summary: "Create log for asset",
            path: "/assets/:id/logs",
            method: "POST",
            description: "Create a new log entry for a specific asset",
            pathParams: { id: ZUuid },
            body: ZCreateLogRequest,
            responses: {
                201: ZAssetLog,
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
            pathParams: { id: ZUuid },
            responses: {
                200: ZAssetLog,
                404: ZErrorResponse,
            },
            metadata: metadata,
        },

        updateLog: {
            summary: "Update log",
            path: "/logs/:id",
            method: "PATCH",
            description: "Update an existing log entry (partial update)",
            pathParams: { id: ZUuid },
            body: ZUpdateLogRequest,
            responses: {
                200: ZAssetLog,
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
            pathParams: { id: ZUuid },
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
