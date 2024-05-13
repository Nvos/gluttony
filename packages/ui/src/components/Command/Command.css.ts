import { globalStyle, style } from '@vanilla-extract/css';
import { atom, vars } from '../../theme';

export const root = style([
  atom({
    display: 'flex',
    width: 'full',
    height: 'full',
    flexDirection: 'column',
    overflow: 'hidden',
    borderRadius: 100,
  }),
  {
    backgroundColor: vars.color.surface[50],
  },
]);

export const input = style({
  selectors: {
    '&:placeholder': {
      fontSize: vars.fontSize[100],
    },
  },
});

export const item = style([
  atom({
    position: 'relative',
    display: 'flex',
    alignItems: 'center',
    borderRadius: 100,
    paddingX: 100,
    paddingY: 100,
  }),
  {
    selectors: {
      '&[aria-selected="true"]': {
        backgroundColor: vars.color.neutral[200],
      },
    },
  },
]);

export const itemList = style({
  overflowY: 'auto',
  overflowX: 'hidden',
  maxHeight: 300,
});

export const group = style({
  overflow: 'hidden',
  padding: vars.space[50],
  selectors: {
    '&[hidden=""]': {
      display: 'none',
    },
  },
});

globalStyle(`${group} [cmdk-group-heading=""]`, {
  color: vars.color.neutral[800],
  fontSize: vars.fontSize[50],
  letterSpacing: vars.letterSpacing[50],
  padding: vars.space[100],
});

export const shortcut = style({
  marginLeft: 'auto',
  color: vars.color.neutral[800],
});
