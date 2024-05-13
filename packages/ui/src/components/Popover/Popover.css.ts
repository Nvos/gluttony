import { style } from '@vanilla-extract/css';
import { atom, vars } from '../../theme';

export const root = style([
  {
    zIndex: vars.zIndex.modal,
    borderRadius: vars.radii[100],
    padding: vars.space[300],
    boxShadow: vars.shadow[100],
    backgroundColor: vars.color.surface[50],
  },
  atom({ border: 'neutral' }),
]);
