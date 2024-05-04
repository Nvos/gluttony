import { style } from '@vanilla-extract/css';
import { atom, vars } from '../../theme';

export const root = style([
  atom({ position: 'fixed' }),
  {
    zIndex: vars.zIndex.modal,
    inset: 0,
    backgroundColor: 'rgba(0,0,0,.8)',
  },
]);

export const content = style([
  atom({
    position: 'fixed',
    width: 'full',
    borderRadius: 100,
    padding: 300,
  }),
  {
    left: '50%',
    top: '50%',
    zIndex: vars.zIndex.modal,
    boxShadow: vars.shadow[200],
    backgroundColor: vars.color.surface[200],
    transform: 'translate(-50%, -50%)',
    maxWidth: '512px',
  },
]);

export const closeIcon = style([
  atom({
    position: 'absolute',
  }),
  {
    top: vars.space[300],
    right: vars.space[300],
  },
]);
