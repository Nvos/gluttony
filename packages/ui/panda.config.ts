import { defineConfig } from '@pandacss/dev';
import pandaBasePreset from '@pandacss/preset-base';
import { preset } from 'src/preset';

const patterns = pandaBasePreset.patterns;
export default defineConfig({
  presets: ['@pandacss/preset-base', preset],
  preflight: true,
  jsxFramework: 'react',
  include: ['./src/**/*.{js,jsx,ts,tsx}'],
  exclude: [],
  strictTokens: true,
  strictPropertyValues: true,
  outdir: '../theme',
  importMap: '@gluttony/theme',
});
