import { recipe } from '@vanilla-extract/recipes';
import { vars } from './theme.contract.css';

export const typography = recipe({
  variants: {
    color: {
      standard: {
        color: vars.color.neutral[1000],
      },
      destructive: {
        color: vars.color.danger[900],
      },
      caption: {
        color: vars.color.neutral[800],
      },
    },
    weight: {
      normal: { fontWeight: vars.fontWeight.normal },
      medium: { fontWeight: vars.fontWeight.medium },
      strong: { fontWeight: vars.fontWeight.strong },
    },
    align: {
      center: {
        textAlign: 'center',
      },
    },
    size: {
      50: {
        fontSize: vars.fontSize[50],
        letterSpacing: vars.letterSpacing[50],
        lineHeight: '17px',
      },
      100: {
        fontSize: vars.fontSize[100],
        letterSpacing: vars.letterSpacing[100],
        lineHeight: '20px',
      },
      200: {
        fontSize: vars.fontSize[200],
        letterSpacing: vars.letterSpacing[200],
        lineHeight: '21px',
      },
      300: {
        fontSize: vars.fontSize[300],
        letterSpacing: vars.letterSpacing[300],
        lineHeight: '29px',
      },
    },
  },
});
