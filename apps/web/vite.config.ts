import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react'; // TODO: change to swc
import tsconfigPaths from 'vite-tsconfig-paths';
import { TanStackRouterVite } from '@tanstack/router-vite-plugin';

const config = defineConfig({
  plugins: [tsconfigPaths(), react(), TanStackRouterVite({ })],
  server: {
    port: 3000,
  },
});

export default config;
