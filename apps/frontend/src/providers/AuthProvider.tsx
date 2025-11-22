import { useEffect, type ReactNode } from "react";
import { useAuth } from "@clerk/clerk-react";

/**
 * AuthProvider component that ensures Clerk authentication is initialized
 * Token injection is handled automatically by useApiClient hook
 */
export function AuthProvider({ children }: { children: ReactNode }) {
    const { isSignedIn, getToken } = useAuth();

    useEffect(() => {
        // When user is signed in, ensure token is available
        // This is primarily for debugging and ensuring Clerk is ready
        if (isSignedIn) {
            getToken({ template: "custom" }).catch((error) => {
                console.error("Failed to get auth token:", error);
            });
        }
    }, [isSignedIn, getToken]);

    return <>{children}</>;
}
