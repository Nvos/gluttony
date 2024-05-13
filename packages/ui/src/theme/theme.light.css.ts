import { createTheme } from '@vanilla-extract/css';
import { vars } from './theme.contract.css';
import { baseTheme } from './base.css';
import { radixToLocalTokenScale } from './color.util';
import { orange, slate, red } from '@radix-ui/colors';

export const lightTheme = createTheme(vars, {
  ...baseTheme,
  color: {
    primary: radixToLocalTokenScale(orange),
    neutral: radixToLocalTokenScale(slate),
    danger: radixToLocalTokenScale(red),
    surface: {
      50: vars.color.neutral[50],
      100: vars.color.neutral[100],
      200: vars.color.neutral[200],
      300: vars.color.neutral[300],
    },
    text: {
      standard: vars.color.neutral[1000],
      muted: vars.color.neutral[800],
      danger: vars.color.danger[900],
    },
  },
});
