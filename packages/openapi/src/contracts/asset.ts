import { getSecurityMetadata } from "../utils.js";
import {
    ZAsset,
    ZCreateAssetRequest,
    ZUpdateAssetRequest,
    ZAssetListResponse,
    ZAssetQueryParams,
    ZErrorResponse,
    ZUuid,
} from "@gardenjournal/zod";
import { initContract } from "@ts-rest/core";
import { z } from "zod";

const c = initContract();

const metadata = getSecurityMetadata();

export const assetContract = c.router(
    {
        listAssets: {
            summary: "List assets",
            path: "/assets",
            method: "GET",
            description: "Get a paginated list of assets for the authenticated user",
            query: ZAssetQueryParams,
            responses: {
                200: ZAssetListResponse,
            },
            metadata: metadata,
        },

        createAsset: {
            summary: "Create asset",
            path: "/assets",
            method: "POST",
            description: "Create a new asset for the authenticated user",
            body: ZCreateAssetRequest,
            responses: {
                201: ZAsset,
                400: ZErrorResponse,
            },
            metadata: metadata,
        },

        getAssetById: {
            summary: "Get asset by ID",
            path: "/assets/:id",
            method: "GET",
            description: "Get a single asset by its ID",
            pathParams: z.object({ id: ZUuid }),
            responses: {
                200: ZAsset,
                404: ZErrorResponse,
            },
            metadata: metadata,
        },

        updateAsset: {
            summary: "Update asset",
            path: "/assets/:id",
            method: "PATCH",
            description: "Update an existing asset (partial update)",
            pathParams: z.object({ id: ZUuid }),
            body: ZUpdateAssetRequest,
            responses: {
                200: ZAsset,
                400: ZErrorResponse,
                404: ZErrorResponse,
            },
            metadata: metadata,
        },

        deleteAsset: {
            summary: "Delete asset",
            path: "/assets/:id",
            method: "DELETE",
            description: "Delete an asset and all its associated logs",
            pathParams: z.object({ id: ZUuid }),
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
