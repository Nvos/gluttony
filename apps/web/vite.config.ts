import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';
import tsconfigPaths from 'vite-tsconfig-paths';
import { TanStackRouterVite } from '@tanstack/router-vite-plugin';
import { vanillaExtractPlugin } from '@vanilla-extract/vite-plugin';

const config = defineConfig({
  plugins: [
    tsconfigPaths(),
    vanillaExtractPlugin(),
    react(),
    TanStackRouterVite({ enableRouteGeneration: false }),
  ],
  server: {
    port: 3000,
  },
});

export default config;
