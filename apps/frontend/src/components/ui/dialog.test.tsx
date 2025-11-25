import { describe, it, expect } from "vitest";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
    DialogFooter,
} from "./dialog";

describe("Dialog Components", () => {
    it("should export all dialog components", () => {
        expect(Dialog).toBeDefined();
        expect(DialogContent).toBeDefined();
        expect(DialogDescription).toBeDefined();
        expect(DialogHeader).toBeDefined();
        expect(DialogTitle).toBeDefined();
        expect(DialogTrigger).toBeDefined();
        expect(DialogFooter).toBeDefined();
    });

    it("should be functions", () => {
        expect(typeof Dialog).toBe("function");
        expect(typeof DialogContent).toBe("function");
        expect(typeof DialogTitle).toBe("function");
    });
});
