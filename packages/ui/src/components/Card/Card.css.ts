import { style } from '@vanilla-extract/css';
import { atom, vars } from '../../theme';

export const root = style([
  atom({
    borderRadius: 100,
  }),
  {
    backgroundColor: vars.color.surface[100],
    border: `solid 1px ${vars.color.neutral[600]}`,
    boxShadow: vars.shadow[100],
  },
]);

export const image = style({
  display: 'block',
  backgroundSize: 'cover',
  backgroundRepeat: 'no-repeat',
  backgroundPosition: 'center',
  width: 'full',
  objectFit: 'cover',
});
