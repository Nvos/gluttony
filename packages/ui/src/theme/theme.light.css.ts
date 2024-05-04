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
      50: '',
      100: '',
      200: '',
      300: '',
    },
  },
});
