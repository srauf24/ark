import "./index.css";
import {
  ClerkProvider,
  SignIn,
  SignedIn,
  SignedOut,
  RedirectToSignIn,
} from "@clerk/clerk-react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AuthProvider } from "@/providers/AuthProvider";
import { CLERK_PUBLISHABLE_KEY } from "@/config/env";

// Create QueryClient instance with default options
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 1000 * 60 * 5, // 5 minutes
    },
  },
});

function App() {
  return (
    <ClerkProvider publishableKey={CLERK_PUBLISHABLE_KEY}>
      <QueryClientProvider client={queryClient}>
        <AuthProvider>
          <BrowserRouter>
            <Routes>
              {/* Sign-in route */}
              <Route path="/sign-in/*" element={<SignIn />} />

              {/* Protected routes */}
              <Route
                path="/*"
                element={
                  <>
                    <SignedIn>
                      <div className="p-8">
                        <h1 className="text-3xl font-bold">Authenticated!</h1>
                        <p className="mt-4 text-gray-600">
                          You are now signed in to Ark.
                        </p>
                      </div>
                    </SignedIn>
                    <SignedOut>
                      <RedirectToSignIn />
                    </SignedOut>
                  </>
                }
              />
            </Routes>
          </BrowserRouter>
        </AuthProvider>
      </QueryClientProvider>
    </ClerkProvider>
  );
}

export default App;
