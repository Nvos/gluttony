import { vars } from '@gluttony/design-system';
import { style } from '@vanilla-extract/css';

export const root = style({
  width: '320px',
  minHeight: '100dvh',
  maxHeight: '100dvh',
  display: 'flex',
  flexDirection: 'column',
  gap: '16px',
  padding: '16px',
});

export const section = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '8px',
});

export const link = style({
  textDecoration: 'none',
  padding: `8px 4px`,
  borderRadius: 5,
  color: vars.color.background[800],
  selectors: {
    '&.active': {
      color: vars.color.text[600],
      backgroundColor: vars.color.primary[600],
    },

    '&.active:hover': {
      transition: '250ms',
      color: vars.color.text[700],
      backgroundColor: vars.color.primary[700],
    },
  },
});
