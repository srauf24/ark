import { describe, it, expect } from "vitest";
import { Toaster } from "./sonner";

describe("Toaster Component", () => {
    it("should export Toaster component", () => {
        expect(Toaster).toBeDefined();
    });

    it("should be a function", () => {
        expect(typeof Toaster).toBe("function");
    });
});
