import tailwindcss from "@tailwindcss/vite";
import path from "path";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { router } from "sv-router/vite-plugin";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [tailwindcss(), svelte(), router()],
  resolve: {
    alias: {
      $lib: path.resolve("./src/lib"),
      $src: path.resolve("./src"),
    },
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
});
