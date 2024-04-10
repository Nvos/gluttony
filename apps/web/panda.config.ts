import { defineConfig } from '@pandacss/dev';

export default defineConfig({
  presets: ['@pandacss/preset-base', '@gluttony/ui/preset'],
  preflight: true,
  jsxFramework: 'react',
  exclude: [],
  include: ['./src/**/*.{ts,tsx,js,jsx}', '../../packages/ui/src/**/*.tsx'],
  outdir: '../../packages/theme',
  importMap: '@glutony/theme',
});
