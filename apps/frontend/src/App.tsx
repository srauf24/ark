import "./index.css";
import {
  ClerkProvider,
  SignIn,
  SignedIn,
  SignedOut,
} from "@clerk/clerk-react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "@/providers/AuthProvider";
import { CLERK_PUBLISHABLE_KEY } from "@/config/env";
import { AssetList } from "@/components/assets/AssetList";
import { Layout } from "@/components/layout/Layout";
import { AssetDetailPage } from "@/pages/assets/AssetDetailPage";
import { LandingPage } from "@/pages/LandingPage";
import { Toaster } from "@/components/ui/sonner";

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
              <Route path="/sign-in/*" element={<SignIn routing="path" path="/sign-in" />} />
              <Route path="/sign-up/*" element={<SignIn routing="path" path="/sign-up" />} />

              {/* Landing page for signed-out users, redirect to /assets for signed-in */}
              <Route
                path="/"
                element={
                  <>
                    <SignedOut>
                      <LandingPage />
                    </SignedOut>
                    <SignedIn>
                      <Navigate to="/assets" replace />
                    </SignedIn>
                  </>
                }
              />

              {/* Protected routes - Layout with nested routes */}
              <Route element={<Layout />}>
                <Route
                  path="/assets"
                  element={
                    <SignedIn>
                      <AssetList />
                    </SignedIn>
                  }
                />
                <Route
                  path="/assets/:id"
                  element={
                    <SignedIn>
                      <AssetDetailPage />
                    </SignedIn>
                  }
                />
              </Route>
            </Routes>
          </BrowserRouter>
        </AuthProvider>
      </QueryClientProvider>
      <Toaster />
    </ClerkProvider>
  );
}

export default App;
