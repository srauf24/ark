import { describe, it, expect } from "vitest";
import { CLERK_PUBLISHABLE_KEY, API_URL, ENV } from "./env";

describe("Environment Configuration", () => {
    it("should export CLERK_PUBLISHABLE_KEY", () => {
        expect(CLERK_PUBLISHABLE_KEY).toBeDefined();
        expect(typeof CLERK_PUBLISHABLE_KEY).toBe("string");
    });

    it("should export API_URL", () => {
        expect(API_URL).toBeDefined();
        expect(typeof API_URL).toBe("string");
        // Should be a valid URL
        expect(() => new URL(API_URL)).not.toThrow();
    });

    it("should export ENV", () => {
        expect(ENV).toBeDefined();
        expect(["production", "development", "local"]).toContain(ENV);
    });

    it("should have valid API_URL format", () => {
        const url = new URL(API_URL);
        expect(url.protocol).toMatch(/^https?:$/);
    });

    it("should use test defaults in test environment", () => {
        // In test environment, env.ts provides mock values
        expect(CLERK_PUBLISHABLE_KEY).toBeTruthy();
        expect(API_URL).toBeTruthy();
        expect(ENV).toBe("local");
    });
});
