import { recipe } from '@vanilla-extract/recipes';
import { atom, typography, vars } from '../../theme';

export const root = recipe({
  base: [
    atom({
      borderRadius: 100,
      paddingX: 200,
    }),
    typography({ size: 100, color: 'standard' }),
    {
      minWidth: 0,
      outlineOffset: '2px',
      borderWidth: '1px',
      borderStyle: 'solid',
      borderColor: vars.color.neutral[500],
      selectors: {
        '&:focus-visible': {
          outline: `2px solid ${vars.color.primary[800]}`,
        },
      },
    },
  ],
  variants: {
    size: {
      sm: {
        height: vars.size[100],
      },
      md: {
        height: vars.size[200],
      },
    },
  },
  defaultVariants: {
    size: 'md',
  },
});
