import { TextStyles } from '@pandacss/dev';

export const fontWeights = {
  normal: { value: 400 },
  medium: { value: 500 },
  heavy: { value: 700 },
};

export const textStyles: TextStyles = {
  sm: {
    value: {
      fontSize: '0.875rem',
      lineHeight: '1.25rem',
      letterSpacing: '-0.006em',
    },
  },
  md: {
    value: {
      fontSize: '1rem',
      lineHeight: '1.5rem',
      letterSpacing: '-0.011em',
    },
  },
  lg: {
    value: {
      fontSize: '1.125rem',
      lineHeight: '1.75rem',
      letterSpacing: '-0.019em',
    },
  },
  xl: {
    value: {
      fontSize: '1.25rem',
      lineHeight: '1.75rem',
      letterSpacing: '-0.021em',
    },
  },
  '2xl': {
    value: {
      fontSize: '1.5rem',
      lineHeight: '2rem',
      letterSpacing: '-0.022em',
    },
    '3xl': {
      value: {
        fontSize: '1.875rem',
        lineHeight: '2.25rem',
        letterSpacing: '-0.022em',
      },
    },
    '4xl': {
      value: {
        fontSize: '2.25rem',
        lineHeight: '2.5rem',
        letterSpacing: '-0.022em',
      },
    },
    '5xl': {
      value: {
        fontSize: '3rem',
        lineHeight: '1',
        letterSpacing: '-0.022em',
      },
    },
  },
};
