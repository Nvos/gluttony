import { defineRecipe } from '@pandacss/dev';

export const buttonRecipe = defineRecipe({
  className: 'button',
  base: {
    display: 'inline-flex',
    alignItems: 'center',
    justifyContent: 'center',
    whiteSpace: 'nowrap',
    borderRadius: '100',
    paddingX: '200',
    textStyle: 'md',
    fontWeight: 'medium',
    userSelect: 'none',
    minWidth: 0,
    '&:hover': {
      transition: '250ms',
      cursor: 'pointer',
    },
  },
  variants: {
    size: {
      sm: { height: '36px' },
      md: {
        height: '40px',
      },
    },
    variant: {
      solid: {
        background: 'colorPalette.800',
        color: 'colorPalette.50',
        '&:hover': {
          background: 'colorPalette.900',
        },
      },
      ghost: {
        background: 'colorPalette.800/20',
        color: 'colorPalette.950',
        '&:hover': {
          background: 'colorPalette.800/30',
        },
      },
    },
  },
  defaultVariants: {
    size: 'md',
    variant: 'solid',
  },
});
