import { definePreset } from '@pandacss/dev';
import { globalCss } from './theme/global';
import { orange, orangeDark, slate, slateDark, red, redDark } from '@radix-ui/colors';
import { radii } from './theme/border';
import { radixToLocalTokenScale, colorTokenToSemanticToken } from './theme/color/converter';
import { buttonRecipe } from './theme/recipe/button';
import { spacing } from './theme/spacing';
import { fontWeights, textStyles } from './theme/typography';
import { inputRecipe } from './theme/recipe/input';

export const preset = definePreset({
  conditions: {
    dark: '[data-theme=dark] &',
  },
  theme: {
    extend: {
      tokens: {
        sizes: {
          full: { value: '100%' },
          min: { value: 'min-content' },
          max: { value: 'max-content' },
          fit: { value: 'fit-content' },
          dvh: { value: '100dvh' },
        },
        spacing: spacing,
        radii: radii,
        fontWeights: fontWeights,
        colors: {
          white: { DEFAULT: { value: '#ffffff' } },
          black: {
            DEFAULT: { value: '#000000' },
          },
          orange: radixToLocalTokenScale(orange, orangeDark),
          slate: radixToLocalTokenScale(slate, slateDark),
          red: radixToLocalTokenScale(red, redDark),
        },
      },
      recipes: {
        button: buttonRecipe,
        input: inputRecipe,
      },
      textStyles: textStyles,
      semanticTokens: {
        shadows: {
          100: {
            value: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
          },
        },
        colors: {
          // 50, 100 <-backgrounds
          // 200, 300, 400 <-interactive
          // 500, 600, 700 <- borders and separators
          // 800, 900 <- solid color
          // 950, 100 <- accessible text
          danger: colorTokenToSemanticToken('red'),
          neutral: colorTokenToSemanticToken('slate'),
          primary: colorTokenToSemanticToken('orange'),
          background: {
            DEFAULT: { value: { base: '{colors.neutral.50}', _dark: '{colors.neutral.50}' } },
            layer: { value: { base: '{colors.white}', _dark: '{colors.slate.dark.50}' } },
          },
        },
      },
    },
  },
  globalCss: globalCss,
  staticCss: {
    css: [
      {
        properties: { colorPalette: ['neutral', 'primary', 'danger'] },
      },
    ],
  },
});

export default preset;
