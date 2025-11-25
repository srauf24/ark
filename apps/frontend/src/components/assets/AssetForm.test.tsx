import { describe, it, expect } from "vitest";
import { AssetForm } from "./AssetForm";

describe("AssetForm", () => {
    it("should export AssetForm component", () => {
        expect(AssetForm).toBeDefined();
    });

    it("should be a function", () => {
        expect(typeof AssetForm).toBe("function");
    });
});
