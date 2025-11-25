import { toast } from "sonner";

/**
 * Error handling utilities for consistent user-facing error messages
 */

interface ApiErrorResponse {
    status: number;
    body?: {
        code?: string;
        message?: string;
        errors?: Array<{ field: string; error: string }>;
        details?: Record<string, string[]>; // Keep for backward compatibility if needed
    };
}

/**
 * Get user-friendly error message based on HTTP status code
 */
export function getErrorMessage(error: unknown): string {
    // Handle API error responses
    if (error && typeof error === "object" && "status" in error) {
        const apiError = error as ApiErrorResponse;

        // Use backend error message if available
        if (apiError.body?.message) {
            return apiError.body.message;
        }

        switch (apiError.status) {
            case 400:
                return "Invalid input. Please check your entries.";
            case 401:
                return "Session expired. Please sign in again.";
            case 404:
                return "Asset not found.";
            case 500:
                return "Something went wrong. Please try again.";
            default:
                return "An unexpected error occurred.";
        }
    }

    // Handle Error instances
    if (error instanceof Error) {
        return error.message;
    }

    // Fallback for unknown errors
    return "An unexpected error occurred.";
}

/**
 * Check if error should trigger sign-out (401 Unauthorized)
 */
export function shouldSignOut(error: unknown): boolean {
    if (error && typeof error === "object" && "status" in error) {
        const apiError = error as ApiErrorResponse;
        return apiError.status === 401;
    }
    return false;
}

/**
 * Extract field-level validation errors from 400 Bad Request responses
 * Returns null if no field errors are present
 */
export function getFieldErrors(
    error: unknown
): Record<string, string> | null {
    if (error && typeof error === "object" && "status" in error) {
        const apiError = error as ApiErrorResponse;

        if (apiError.status === 400 && apiError.body) {
            console.log("[v1.1] Parsing field errors from body:", apiError.body);
            const fieldErrors: Record<string, string> = {};

            // Handle backend "errors" array format
            if (apiError.body.errors && Array.isArray(apiError.body.errors)) {
                apiError.body.errors.forEach((err) => {
                    fieldErrors[err.field] = err.error;
                });
            }
            // Handle legacy "details" object format
            else if (apiError.body.details) {
                for (const [field, messages] of Object.entries(apiError.body.details)) {
                    if (Array.isArray(messages) && messages.length > 0) {
                        fieldErrors[field] = messages[0];
                    }
                }
            }

            return Object.keys(fieldErrors).length > 0 ? fieldErrors : null;
        }
    }

    return null;
}

/**
 * Show error toast with appropriate message
 */
export function showErrorToast(error: unknown): void {
    const message = getErrorMessage(error);
    toast.error(message);
}

/**
 * Show success toast
 */
export function showSuccessToast(message: string): void {
    toast.success(message);
}
