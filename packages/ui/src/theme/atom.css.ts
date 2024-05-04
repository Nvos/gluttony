import { createSprinkles, defineProperties } from '@vanilla-extract/sprinkles';
import { vars } from './theme.contract.css';

const properties = defineProperties({
  properties: {
    display: ['flex', 'block', 'inline-flex', 'none'],
    justifyContent: ['center', 'space-between'],
    flexDirection: ['row', 'column'],
    flexWrap: ['wrap'],
    alignItems: ['center'],
    whiteSpace: ['nowrap'],
    position: ['absolute', 'relative'],
    height: {
      0: '0px',
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
