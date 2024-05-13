import { createThemeContract } from '@vanilla-extract/css';

export const vars = createThemeContract({
  zIndex: {
    modal: '',
  },
  space: {
    50: '',
    100: '',
    200: '',
    300: '',
    400: '',
    500: '',
    600: '',
    700: '',
  },
  size: {
    25: '',
    100: '',
    200: '',
  },
  fontSize: {
    50: '',
    100: '',
    200: '',
    300: '',
  },
  letterSpacing: {
    50: '',
    100: '',
    200: '',
    300: '',
  },
  fontFamily: {
    main: '',
  },
  fontWeight: {
    normal: '',
    medium: '',
    strong: '',
  },
  radii: {
    100: '',
  },
  shadow: {
    50: '0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)',
    100: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
    200: '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
  },
  color: {
    primary: {
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
      1000: '',
    },
    danger: {
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
      1000: '',
    },
    neutral: {
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
      1000: '',
    },
    surface: {
      50: '',
      100: '',
      200: '',
      300: '',
    },
    text: {
      standard: '',
      muted: '',
      danger: ''
    }
  },
});
