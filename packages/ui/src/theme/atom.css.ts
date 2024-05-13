import { createSprinkles, defineProperties } from '@vanilla-extract/sprinkles';
import { vars } from './theme.contract.css';

const properties = defineProperties({
  properties: {
    boxShadow: vars.shadow,
    display: ['flex', 'block', 'inline-flex', 'none'],
    justifyContent: ['center', 'space-between', 'end'],
    flexDirection: ['row', 'column'],
    flexWrap: ['wrap'],
    alignItems: ['center'],
    whiteSpace: ['nowrap'],
    position: ['absolute', 'relative', 'fixed'],
    textAlign: ['center', 'left', 'right'],
    overflow: ['hidden', 'auto'],
    border: {
      neutral: `solid 1px ${vars.color.neutral[600]}`,
    },
    borderBottom: {
      neutral: `solid 1px ${vars.color.neutral[600]}`,
    },
    height: {
      0: '0px',
      25: vars.size[25],
      100: vars.size[100],
      200: vars.size[200],
      dvh: '100dvh',
      full: '100%',
    },
    minHeight: {
      dvh: '100dvh',
    },
    width: {
      full: '100%',
      0: '0px',
      25: vars.size[25],
      100: vars.size[100],
      200: vars.size[200],
    },
    borderRadius: {
      100: vars.radii[100],
    },
    flex: {
      1: 1,
    },
    fontWeight: vars.fontWeight,
    paddingTop: vars.space,
    paddingBottom: vars.space,
    paddingLeft: vars.space,
    paddingRight: vars.space,

    marginTop: vars.space,
    marginBottom: vars.space,
    marginLeft: vars.space,
    marginRight: vars.space,

    gap: vars.space,
  },
  shorthands: {
    padding: ['paddingTop', 'paddingBottom', 'paddingLeft', 'paddingRight'],
    paddingX: ['paddingLeft', 'paddingRight'],
    paddingY: ['paddingTop', 'paddingBottom'],

    margin: ['marginTop', 'marginBottom', 'marginLeft', 'marginRight'],
    marginX: ['marginLeft', 'marginRight'],
    marginY: ['marginTop', 'marginBottom'],
  },
});

export const atom = createSprinkles(properties);
