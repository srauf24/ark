import { z } from "zod";

const envVarsSchema = z.object({
  VITE_CLERK_PUBLISHABLE_KEY: z
    .string()
    .min(1, "VITE_CLERK_PUBLISHABLE_KEY is required"),
  VITE_API_URL: z.string().url().default("http://localhost:3000"),
  VITE_ENV: z.enum(["production", "development", "local"]).default("local"),
});

// In test environment, use mock values
const isTest = import.meta.env.MODE === "test";

let envVars: z.infer<typeof envVarsSchema>;

if (isTest) {
  // Provide test defaults
  envVars = {
    VITE_CLERK_PUBLISHABLE_KEY: import.meta.env.VITE_CLERK_PUBLISHABLE_KEY || "pk_test_mock_key",
    VITE_API_URL: import.meta.env.VITE_API_URL || "http://localhost:8080",
    VITE_ENV: "local" as const,
  };
} else {
  const parseResult = envVarsSchema.safeParse(import.meta.env);

  if (!parseResult.success) {
    console.error(
      "❌ Invalid environment variables:",
      z.treeifyError(parseResult.error),
    );
    throw new Error("Invalid environment variables");
  }

  envVars = parseResult.data;
}

// export individual variables
export const ENV = envVars.VITE_ENV;
export const API_URL = envVars.VITE_API_URL;
export const CLERK_PUBLISHABLE_KEY = envVars.VITE_CLERK_PUBLISHABLE_KEY;