import { globalStyle } from '@vanilla-extract/css';
import { vars } from './contract.css';

globalStyle('body', {
  margin: 0,
  fontFamily: `'Roboto', sans-serif`,
  backgroundColor: vars.color.background['50'],
  color: vars.color.background[950],
  fontSize: '16px',
});

globalStyle('*', {
  boxSizing: 'border-box',
});
