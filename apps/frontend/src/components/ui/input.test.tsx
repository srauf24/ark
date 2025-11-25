import { describe, it, expect } from "vitest";
import { Input } from "./input";

describe("Input Component", () => {
    it("should export Input component", () => {
        expect(Input).toBeDefined();
    });

    it("should be a function", () => {
        expect(typeof Input).toBe("function");
    });
});
