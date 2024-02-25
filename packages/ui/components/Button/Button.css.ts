import { style } from '@vanilla-extract/css';
import { vars } from '@gluttony/design-system';

export const root = style({
  height: '40px',
  padding: `${vars.space[2]}`,
  display: 'flex',
  alignItems: 'center',
});
