import { test, expect } from "@playwright/test";

test.describe("Claude AI Navigation", () => {
  test("should navigate to claude.ai", async ({ page }) => {
    // Navigate to Claude AI
    await page.goto("https://claude.ai");

    // Wait for the page to load
    await page.waitForLoadState("domcontentloaded");

    // Verify we're on the Claude AI site
    await expect(page).toHaveURL(/claude\.ai/);

    // Take a screenshot for verification
    await page.screenshot({ path: "claude-ai-screenshot.png", fullPage: true });

    console.log("Successfully navigated to claude.ai");
  });
});
