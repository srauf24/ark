import { extendZodWithOpenApi } from "@anatine/zod-openapi";
import { z } from "zod";

// Only extend Zod with OpenAPI in Node.js environment (for spec generation)
// In browser, this will fail silently as it's not needed for runtime validation
try {
    extendZodWithOpenApi(z);
} catch (e) {
    // Ignore errors in browser environment
}

export * from "./utils.js";
export * from "./common.js";
export * from "./health.js";
export * from "./asset.js";
export * from "./log.js";