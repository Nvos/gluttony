import { globalStyle } from '@vanilla-extract/css';
import { reset } from './layers.css';
import { vars } from './theme.contract.css';

globalStyle(
  '*:where(:not(html, iframe, canvas, img, svg, video, audio):not(svg *, symbol *))',
  {
    '@layer': {
      [reset]: {
        all: 'unset',
        display: 'revert',
      },
    },
  },
);

globalStyle('*, *::before, *::after', {
  '@layer': {
    [reset]: {
      boxSizing: 'border-box',
    },
  },
});

globalStyle('html, body', {
  '@layer': {
    [reset]: {
      MozTextSizeAdjust: 'none',
      WebkitTextSizeAdjust: 'none',
      textSizeAdjust: 'none',

      fontFamily: vars.fontFamily.main,
      fontFeatureSettings: "'liga' 1, 'calt' 1",
      color: vars.color.neutral[1000],
      backgroundColor: vars.color.surface[50],
    },
  },
});

globalStyle('a, button', {
  '@layer': {
    [reset]: {
      cursor: 'pointer',
    },
  },
});

globalStyle('nav, ol, ul, menu, summary', {
  '@layer': {
    [reset]: {
      listStyle: 'none',
    },
  },
});

globalStyle('img', {
  '@layer': {
    [reset]: {
      maxInlineSize: '100%',
      maxBlockSize: '100%',
    },
  },
});

globalStyle('table', {
  '@layer': {
    [reset]: {
      borderCollapse: 'collapse',
    },
  },
});
