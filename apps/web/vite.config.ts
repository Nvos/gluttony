import { defineConfig } from '@solidjs/start/config';
import { vanillaExtractPlugin } from '@vanilla-extract/vite-plugin';

export default defineConfig({
  server: {
    fs: {
      allow: [
        'E:/dev/gluttony/node_modules/.pnpm/@solidjs+start@0.5.9_solid-js@1.8.15_vinxi@0.2.1_vite@5.1.3/node_modules/@solidjs/start/shared/dev-overlay/DevOverlayDialog.tsx',
      ],
    },
  },
  plugins: [vanillaExtractPlugin()],
  start: {
    ssr: false,
  },
});
