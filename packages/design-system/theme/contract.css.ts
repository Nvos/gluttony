import { createThemeContract } from '@vanilla-extract/css';

const colorContract = {
  50: '',
  100: '',
  200: '',
  300: '',
  400: '',
  500: '',
  600: '',
  700: '',
  800: '',
  900: '',
  950: '',
};

const colorCommonContract = {
  white: '',
  black: '',
};

const bodyFontSizeContract = {
  '3xl': '',
  '2xl': '',
  xl: '',
  lg: '',
  md: '',
  sm: '',
  xs: '',
};

const headingFontSizeContract = {
  '3xl': '',
  '2xl': '',
  xl: '',
  lg: '',
  md: '',
  sm: '',
  xs: '',
  '2xs': '',
};

// Increments of 4
const spaceContract = {
  0: '',
  1: '', // 4
  2: '', // 8
  4: '', // 16
  8: '', // 32,
  16: '', // 64,
};

const fontWeightContract = {
  regular: '',
  bold: '',
};

const shadowContract = {
  lg: '',
  md: '',
  none: '',
};

export const vars = createThemeContract({
  color: {
    primary: colorContract,
    danger: colorContract,
    warn: colorContract,
    background: colorContract,
    common: colorCommonContract,
    text: colorContract,
  },
  fontSize: {
    body: bodyFontSizeContract,
    heading: headingFontSizeContract,
  },
  space: spaceContract,
  fontWeight: fontWeightContract,
  shadow: shadowContract,
});
