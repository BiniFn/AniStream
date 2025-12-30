import { defineConfig } from "tsup";

const isDev = process.env.NODE_ENV === "development";

export default defineConfig({
  entry: {
    index: "src/index.ts",
    preload: "src/preload.ts",
  },
  format: ["cjs"],
  target: "node20",
  platform: "node",
  outDir: "dist",
  clean: true,
  sourcemap: true,
  shims: true,
  external: ["electron"],
  define: {
    "process.env.ANIWAYS_URL": JSON.stringify(
      process.env.ANIWAYS_URL || (isDev ? "http://localhost:3000" : "https://aniways.xyz")
    ),
    "process.env.NODE_ENV": JSON.stringify(process.env.NODE_ENV || "production"),
  },
});
