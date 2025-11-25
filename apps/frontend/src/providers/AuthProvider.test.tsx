import { describe, it, expect, vi, beforeEach } from "vitest";
import { render } from "@testing-library/react";
import { AuthProvider } from "./AuthProvider";

// Mock Clerk's useAuth hook
const mockGetToken = vi.fn();
const mockUseAuth = vi.fn();

vi.mock("@clerk/clerk-react", () => ({
    useAuth: () => mockUseAuth(),
}));

describe("AuthProvider", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("should render children", () => {
        mockUseAuth.mockReturnValue({
            isSignedIn: false,
            getToken: mockGetToken,
        });

        const { getByText } = render(
            <AuthProvider>
                <div>Test Child</div>
            </AuthProvider>
        );

        expect(getByText("Test Child")).toBeDefined();
    });

    it("should not crash when Clerk is not initialized", () => {
        mockUseAuth.mockReturnValue({
            isSignedIn: false,
            getToken: mockGetToken,
        });

        expect(() => {
            render(
                <AuthProvider>
                    <div>Test</div>
                </AuthProvider>
            );
        }).not.toThrow();
    });

    it("should handle signed-in state", () => {
        mockGetToken.mockResolvedValue("mock-token");
        mockUseAuth.mockReturnValue({
            isSignedIn: true,
            getToken: mockGetToken,
        });

        const { getByText } = render(
            <AuthProvider>
                <div>Authenticated Content</div>
            </AuthProvider>
        );

        expect(getByText("Authenticated Content")).toBeDefined();
    });

    it("should handle signed-out state", () => {
        mockUseAuth.mockReturnValue({
            isSignedIn: false,
            getToken: mockGetToken,
        });

        const { getByText } = render(
            <AuthProvider>
                <div>Public Content</div>
            </AuthProvider>
        );

        expect(getByText("Public Content")).toBeDefined();
        expect(mockGetToken).not.toHaveBeenCalled();
    });

    it("should call getToken when signed in", async () => {
        mockGetToken.mockResolvedValue("mock-token");
        mockUseAuth.mockReturnValue({
            isSignedIn: true,
            getToken: mockGetToken,
        });

        render(
            <AuthProvider>
                <div>Test</div>
            </AuthProvider>
        );

        // Wait for useEffect to run
        await new Promise((resolve) => setTimeout(resolve, 0));

        expect(mockGetToken).toHaveBeenCalledWith({ template: "custom" });
    });
});
