import { defineRecipe } from '@pandacss/dev';

export const inputRecipe = defineRecipe({
  className: 'input',
  base: {
    display: 'flex',
    borderRadius: '100',
    paddingX: '200',
    textStyle: 'md',
    minWidth: 0,
    borderColor: 'neutral.500',
    borderStyle: 'solid',
    borderWidth: '1px',
    outlineOffset: '50',
    '&:focus-visible': {
      outline: '2px solid {colors.primary.800}',
    },
  },
  variants: {
    size: {
      sm: {
        height: '36px',
      },
      md: {
        height: '40px',
      },
    },
  },
  defaultVariants: {
    size: 'md',
  },
});
