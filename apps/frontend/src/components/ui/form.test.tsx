import { describe, it, expect } from "vitest";
import {
    Form,
    FormControl,
    FormDescription,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "./form";

describe("Form Components", () => {
    it("should export all form components", () => {
        expect(Form).toBeDefined();
        expect(FormControl).toBeDefined();
        expect(FormDescription).toBeDefined();
        expect(FormField).toBeDefined();
        expect(FormItem).toBeDefined();
        expect(FormLabel).toBeDefined();
        expect(FormMessage).toBeDefined();
    });

    it("should be functions or objects", () => {
        expect(typeof Form).toBe("function");
        expect(typeof FormField).toBe("function");
        expect(typeof FormItem).toBe("function");
    });
});
