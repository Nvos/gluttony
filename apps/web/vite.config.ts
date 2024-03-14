import { PluginOption, defineConfig } from 'vite';
import { vanillaExtractPlugin } from '@vanilla-extract/vite-plugin';
import viteSolidPlugin from 'vite-plugin-solid';

const config = defineConfig(({}) => {
  return {
    plugins: [vanillaExtractPlugin(), viteSolidPlugin()],
    server: {
      port: 3000,
    },
    css: {
      transformer: 'lightningcss',
      lightningcss: {
        targets: { firefox: 112, chrome: 120 },
      },
    },
  };
});

export default config;
