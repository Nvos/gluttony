import { recipe, RecipeVariants } from '@vanilla-extract/recipes';
import { atom, typography, vars } from '../../theme';

export const root = recipe({
  base: [
    atom({
      gap: 100,
      paddingX: 200,
      display: 'inline-flex',
      alignItems: 'center',
      justifyContent: 'center',
      borderRadius: 100,
      whiteSpace: 'nowrap',
    }),
    {
      minWidth: 0,
      userSelect: 'none',
      fontSize: vars.fontSize[100],
      fontFamily: vars.fontFamily.main,
      selectors: {
        '&:hover': {
          transition: '250ms',
          cursor: 'pointer',
        },
      },
    },
    typography({ size: 100, weight: 'medium' }),
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

    variant: {
      primary: {
        background: vars.color.primary[800],
        color: vars.color.primary[50],
        '&:hover': {
          background: vars.color.primary[900],
        },
      },
      secondary: {
        background: vars.color.neutral[300],
        color: vars.color.neutral[1000],
        '&:hover': {
          background: vars.color.neutral[400],
        },
      },
      destructive: {
        background: vars.color.danger[800],
        color: vars.color.danger[50],
        '&:hover': {
          background: vars.color.danger[900],
        },
      },
    },
  },
  defaultVariants: {
    size: 'md',
  },
});

export type ButtonVariants = RecipeVariants<typeof root>;
