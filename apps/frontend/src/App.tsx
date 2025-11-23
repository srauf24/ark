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
import { AssetList } from "@/components/assets/AssetList";
import { Layout } from "@/components/layout/Layout";
import { AssetDetailPage } from "@/pages/assets/AssetDetailPage";

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
                      <Routes>
                        <Route element={<Layout />}>
                          <Route path="/" element={<AssetList />} />
                          <Route path="/assets" element={<AssetList />} />
                          <Route path="/assets/:id" element={<AssetDetailPage />} />
                        </Route>
                      </Routes>
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
