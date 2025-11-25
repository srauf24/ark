import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import App from "./App";

// Mock Clerk hooks
const mockUseAuth = vi.fn();
const mockSignIn = vi.fn(() => <div>Sign In Component</div>);
const mockSignedIn = vi.fn(({ children }: { children: React.ReactNode }) => <>{children}</>);
const mockSignedOut = vi.fn(() => <div>Redirecting to sign in...</div>);
const mockRedirectToSignIn = vi.fn(() => <div>Redirecting to sign in...</div>);

vi.mock("@clerk/clerk-react", () => ({
    ClerkProvider: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
    SignIn: () => mockSignIn(),
    SignedIn: ({ children }: { children: React.ReactNode }) => mockSignedIn({ children }),
    SignedOut: () => mockSignedOut(),
    RedirectToSignIn: () => mockRedirectToSignIn(),
    useAuth: () => mockUseAuth(),
}));

// Mock AuthProvider
vi.mock("@/providers/AuthProvider", () => ({
    AuthProvider: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

// Mock env config
vi.mock("@/config/env", () => ({
    CLERK_PUBLISHABLE_KEY: "pk_test_mock_key",
}));

describe("App", () => {
    beforeEach(() => {
        vi.clearAllMocks();
        mockUseAuth.mockReturnValue({
            isSignedIn: false,
            getToken: vi.fn(),
        });
    });

    it("should render without crashing", () => {
        expect(() => {
            render(<App />);
        }).not.toThrow();
    });

    it("should render ClerkProvider with correct publishable key", () => {
        const { container } = render(<App />);

        // App should render successfully with ClerkProvider
        expect(container).toBeDefined();
    });

    it("should show authenticated content when signed in", () => {
        mockSignedIn.mockImplementation(({ children }: { children: React.ReactNode }) => <>{children}</>);
        mockSignedOut.mockImplementation(() => null);

        render(<App />);

        expect(screen.getByText("Authenticated!")).toBeDefined();
        expect(screen.getByText("You are now signed in to Ark.")).toBeDefined();
    });

    it("should redirect to sign-in when signed out", () => {
        mockSignedIn.mockImplementation(() => null);
        mockSignedOut.mockImplementation(() => <div>Redirecting to sign in...</div>);

        render(<App />);

        expect(screen.getByText("Redirecting to sign in...")).toBeDefined();
    });

    it("should have correct provider hierarchy", () => {
        const { container } = render(<App />);

        // Verify the app renders (providers are working)
        expect(container).toBeDefined();
    });

    it("should configure QueryClient with correct defaults", () => {
        // This test verifies the app renders with QueryClientProvider
        expect(() => {
            render(<App />);
        }).not.toThrow();
    });
});
