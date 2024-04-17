import { defineConfig } from '@pandacss/dev';
import {} from './src/theme';

export default defineConfig({
  presets: ['@pandacss/preset-base', './src/theme/preset'],
  preflight: true,
  jsxFramework: 'react',
  include: ['./src/**/*.{ts,tsx,js,jsx}'],
  outdir: './src/ui',
  exclude: ['./src/ui'],
});
