import { defineConfig, loadEnv } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { VitePWA as vitePWA } from "vite-plugin-pwa";
import webfontDownload from "vite-plugin-webfont-dl";
import sveltePreprocess from "svelte-preprocess";

import type * as vite from "vite";
import * as crypto from "crypto";
import * as path from "path";
import * as fs from "fs/promises";
import * as os from "os";

import manifest from "./web/manifest.json";

export default defineConfig({
  plugins: [
    {
      name: "manifest-index",
      transformIndexHtml: {
        enforce: "pre",
        transform(html: string) {
          return html.replace(
            /{{\s*(.+)\s*}}/gi,
            // https://esbuild.github.io/content-types/#direct-eval
            (match, expr) =>
              new Function("manifest", `return ${expr}`)(manifest) || ""
          );
        },
      },
    },
    webfontDownload([
      "https://fonts.googleapis.com/css2?family=Source+Sans+Pro:wght@400;600;700&display=swap",
    ]),
    vitePWA({
      manifest,
      registerType: "autoUpdate",
      // See:
      // https://vite-pwa-org.netlify.app/workbox/generate-sw.html
      // https://vite-pwa-org.netlify.app/workbox/inject-manifest.html
      strategies: "generateSW",
      workbox: {
        globPatterns: ["**/*.{js,css,html,ico,png}"],
      },
      devOptions: {
        enabled: true,
        // This plugin is fucking stupid. There's no way to change the directory
        // where it puts the dev-dist folder, which is dumb. Who the fuck
        // thought putting temp garbage in src was a good idea? Just mktemp!
        resolveTempFolder: () => path.resolve(__dirname, "dist/.pwa-dev"),
      },
    }),
    svelte({
      preprocess: sveltePreprocess(),
    }),
  ],
  worker: {
    format: "iife",
  },
  root: path.resolve(__dirname, "web"),
  envPrefix: "APP_",
  publicDir: "public",
  server: {
    port: 5002,
    host: true,
  },
  build: {
    outDir: "../dist",
    emptyOutDir: true,
  },
  // https://github.com/vitejs/vite/issues/7385#issuecomment-1286606298
  resolve: {
    alias: {
      "#": path.resolve(__dirname, "./web/"),
    },
  },
});
