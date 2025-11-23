import { defineConfig } from "vitest/config";
import path from "path";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 3000,
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@ark/openapi": path.resolve(
        __dirname,
        "../../packages/openapi/dist"
      ),
      "@ark/zod": path.resolve(__dirname, "../../packages/zod/dist"),
    },
  },
  test: {
    globals: true,
    environment: "happy-dom",
    setupFiles: ["./src/test/setup.ts"],
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@ark/openapi": path.resolve(__dirname, "../../packages/openapi/dist"),
      "@ark/zod": path.resolve(__dirname, "../../packages/zod/dist"),
    },
  },
});
