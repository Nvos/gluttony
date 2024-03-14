import { vars } from '@gluttony/design-system';
import { recipe } from '@vanilla-extract/recipes';

export const root = recipe({
  base: {
    height: '40px',
    padding: `${vars.space[2]}`,
    display: 'flex',
    alignItems: 'center',
    border: 'none',
    borderRadius: 5,
    fontSize: 16,
    fontWeight: 700,
    ':focus-visible': {
      outlineOffset: '2px',
      outlineWidth: '2px',
      outlineStyle: 'solid',
      outlineColor: vars.color.primary[500],
      transition: '250ms',
    },
    ':hover': {
      cursor: 'pointer',
    },
  },
  variants: {
    variant: {
      solid: {},
      outline: {},
    },
    colorScheme: {
      background: {},
      primary: {},
      danger: {},
    },
  },
  compoundVariants: [
    // Solid
    {
      variants: { colorScheme: 'background', variant: 'solid' },
      style: {
        backgroundColor: vars.color.background[600],
        color: vars.color.text[600],

        ':hover': {
          backgroundColor: vars.color.background[700],
          transition: '250ms',
        },
      },
    },
    {
      variants: { colorScheme: 'primary', variant: 'solid' },
      style: {
        backgroundColor: vars.color.primary[600],
        color: vars.color.text[600],

        ':hover': {
          backgroundColor: vars.color.primary[700],
          transition: '250ms',
        },
      },
    },
    {
      variants: { colorScheme: 'danger', variant: 'solid' },
      style: {
        backgroundColor: vars.color.danger[600],
        color: vars.color.text[600],

        ':hover': {
          backgroundColor: vars.color.danger[700],
          transition: '250ms',
        },
      },
    },
    // outline
    {
      variants: { colorScheme: 'background', variant: 'outline' },
      style: {
        backgroundColor: 'transparent',
        color: vars.color.background[800],
        border: `solid 2px ${vars.color.background[800]}`,
        ':hover': {
          backgroundColor: vars.color.background[100],
          transition: '250ms',
        },
      },
    },
    {
      variants: { colorScheme: 'primary', variant: 'outline' },
      style: {
        backgroundColor: 'transparent',
        color: vars.color.primary[800],
        border: `solid 2px ${vars.color.primary[800]}`,
        ':hover': {
          backgroundColor: vars.color.primary[100],
          transition: '250ms',
        },
      },
    },
    {
      variants: { colorScheme: 'danger', variant: 'outline' },
      style: {
        backgroundColor: 'transparent',
        color: vars.color.danger[800],
        border: `solid 2px ${vars.color.danger[800]}`,
        ':hover': {
          backgroundColor: vars.color.danger[100],
          transition: '250ms',
        },
      },
    },
  ],
});
