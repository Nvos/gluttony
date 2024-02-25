import { vars } from '@gluttony/design-system';
import { style } from '@vanilla-extract/css';

export const root = style({
  display: 'flex',
  flexDirection: 'column',
  padding: vars.space[4],
  backgroundColor: vars.color.common.white,
  boxShadow: vars.shadow.md,
  borderRadius: 5,
});
