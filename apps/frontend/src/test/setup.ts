// Set environment variables BEFORE any imports
process.env.VITE_CLERK_PUBLISHABLE_KEY = "pk_test_mock_key_for_testing";
process.env.VITE_API_URL = "http://localhost:8080";
process.env.VITE_ENV = "local";

import "@testing-library/jest-dom/vitest";
