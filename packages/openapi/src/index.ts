import { extendZodWithOpenApi } from "@anatine/zod-openapi";
import { z } from "zod";

extendZodWithOpenApi(z);
import { generateOpenApi } from "@ts-rest/open-api";

import { apiContract } from "./contracts/index.js";

type SecurityRequirementObject = {
  [key: string]: string[];
};

export type OperationMapper = NonNullable<
  Parameters<typeof generateOpenApi>[2]
>["operationMapper"];

const hasSecurity = (
  metadata: unknown
): metadata is { openApiSecurity: SecurityRequirementObject[] } => {
  return (
    !!metadata && typeof metadata === "object" && "openApiSecurity" in metadata
  );
};

const operationMapper: OperationMapper = (operation, appRoute) => ({
  ...operation,
  ...(hasSecurity(appRoute.metadata)
    ? {
      security: appRoute.metadata.openApiSecurity,
    }
    : {}),
});

export const OpenAPI = Object.assign(
  generateOpenApi(
    apiContract,
    {
      openapi: "3.0.2",
      info: {
        version: "1.0.0",
        title: "ARK Asset Management API",
        description: "REST API for ARK - A homelab asset tracking and configuration log management application. Track servers, VMs, containers, and network equipment while maintaining searchable logs of configuration changes.",
      },
      servers: [
        {
          url: "http://localhost:8080",
          description: "Local Server",
        },
      ],
    },
    {
      operationMapper,
      setOperationId: true,
    }
  ),
  {
    components: {
      securitySchemes: {
        bearerAuth: {
          type: "http",
          scheme: "bearer",
          bearerFormat: "JWT",
        },
        "x-service-token": {
          type: "apiKey",
          name: "x-service-token",
          in: "header",
        },
      },
    },
  }
);
