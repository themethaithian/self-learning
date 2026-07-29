import { fileURLToPath } from "node:url";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vitest/config";

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL(".", import.meta.url)),
    },
  },
  test: {
    // Plain .test.ts files (lib/*) stay on the faster "node" environment;
    // component tests opt into jsdom per-file via a `@vitest-environment
    // jsdom` docblock instead of paying jsdom's cost repo-wide.
    environment: "node",
    setupFiles: ["./vitest.setup.ts"],
    include: ["**/*.test.ts", "**/*.test.tsx"],
    exclude: ["node_modules/**", ".next/**", "out/**"],
  },
});
