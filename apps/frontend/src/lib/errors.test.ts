import { describe, it, expect } from "vitest";
import {
    getErrorMessage,
    shouldSignOut,
    getFieldErrors,
} from "./errors";

describe("Error Utilities", () => {
    describe("getErrorMessage", () => {
        it("should return correct message for 400 error", () => {
            const error = { status: 400 };
            expect(getErrorMessage(error)).toBe(
                "Invalid input. Please check your entries."
            );
        });

        it("should return correct message for 401 error", () => {
            const error = { status: 401 };
            expect(getErrorMessage(error)).toBe(
                "Session expired. Please sign in again."
            );
        });

        it("should return correct message for 404 error", () => {
            const error = { status: 404 };
            expect(getErrorMessage(error)).toBe("Asset not found.");
        });

        it("should return correct message for 500 error", () => {
            const error = { status: 500 };
            expect(getErrorMessage(error)).toBe(
                "Something went wrong. Please try again."
            );
        });

        it("should return generic message for unknown status code", () => {
            const error = { status: 418 };
            expect(getErrorMessage(error)).toBe("An unexpected error occurred.");
        });

        it("should handle Error instances", () => {
            const error = new Error("Custom error message");
            expect(getErrorMessage(error)).toBe("Custom error message");
        });

        it("should return generic message for unknown error types", () => {
            const error = "string error";
            expect(getErrorMessage(error)).toBe("An unexpected error occurred.");
        });
    });

    describe("shouldSignOut", () => {
        it("should return true for 401 error", () => {
            const error = { status: 401 };
            expect(shouldSignOut(error)).toBe(true);
        });

        it("should return false for non-401 errors", () => {
            expect(shouldSignOut({ status: 400 })).toBe(false);
            expect(shouldSignOut({ status: 404 })).toBe(false);
            expect(shouldSignOut({ status: 500 })).toBe(false);
        });

        it("should return false for non-API errors", () => {
            expect(shouldSignOut(new Error("test"))).toBe(false);
            expect(shouldSignOut("error")).toBe(false);
            expect(shouldSignOut(null)).toBe(false);
        });
    });

    describe("getFieldErrors", () => {
        it("should extract field errors from 400 response", () => {
            const error = {
                status: 400,
                body: {
                    details: {
                        name: ["Name is required"],
                        email: ["Invalid email format"],
                    },
                },
            };

            const fieldErrors = getFieldErrors(error);
            expect(fieldErrors).toEqual({
                name: "Name is required",
                email: "Invalid email format",
            });
        });

        it("should take first error message when multiple exist", () => {
            const error = {
                status: 400,
                body: {
                    details: {
                        password: ["Too short", "Must contain numbers"],
                    },
                },
            };

            const fieldErrors = getFieldErrors(error);
            expect(fieldErrors).toEqual({
                password: "Too short",
            });
        });

        it("should return null for non-400 errors", () => {
            const error = { status: 500 };
            expect(getFieldErrors(error)).toBeNull();
        });

        it("should return null when no details present", () => {
            const error = {
                status: 400,
                body: {},
            };
            expect(getFieldErrors(error)).toBeNull();
        });

        it("should return null for empty details", () => {
            const error = {
                status: 400,
                body: {
                    details: {},
                },
            };
            expect(getFieldErrors(error)).toBeNull();
        });

        it("should return null for non-API errors", () => {
            expect(getFieldErrors(new Error("test"))).toBeNull();
            expect(getFieldErrors("error")).toBeNull();
        });
    });
});
