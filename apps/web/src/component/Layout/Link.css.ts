import { atom, typography, vars } from '@gluttony/ui';
import { style } from '@vanilla-extract/css';

export const linkRoot = style([
  atom({
    display: 'flex',
    gap: 200,
    alignItems: 'center',
    padding: 200,
    height: 200,
    borderRadius: 100,
    width: 'full',
  }),
  typography({ size: 100 }),
  {
    transition: 'background-color 250ms ease',
    selectors: {
      '&.active': {
        color: vars.color.primary[900],
        fontWeight: vars.fontWeight.medium,
      },
      '&.active:hover': {
        backgroundColor: `color-mix(in srgb, ${vars.color.primary[900]}, transparent 95%)`,
      },
      '&:hover': {
        backgroundColor: vars.color.neutral[200],
      },
    },
  },
]);
