import { vars } from './theme.contract.css';

type Scale<T extends string | number> = {
  [P in T]: string;
};

type BaseTheme = {
  space: Scale<keyof typeof vars.space>;
  size: Scale<keyof typeof vars.size>;
  fontSize: Scale<keyof typeof vars.fontSize>;
  letterSpacing: Scale<keyof typeof vars.letterSpacing>;
  fontWeight: Scale<keyof typeof vars.fontWeight>;
  radii: Scale<keyof typeof vars.radii>;
  fontFamily: Scale<keyof typeof vars.fontFamily>;
  shadow: Scale<keyof typeof vars.shadow>;
  zIndex: Scale<keyof typeof vars.zIndex>;
};

export const baseTheme: BaseTheme = {
  zIndex: {
    modal: '1000',
  },
  fontFamily: {
    main: 'Inter, sans-serif',
  },
  shadow: {
    50: '0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)',
    100: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
    200: '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
  },
  space: {
    50: '4px',
    100: '8px',
    200: '16px',
    300: '24px',
    400: '32px',
    500: '48px',
    600: '64px',
    700: '128px',
  },
  size: {
    100: '32px',
    200: '40px',
  },
  fontSize: {
    50: '0.875rem',
    100: '1rem',
    200: '1.125rem',
    300: '1.5rem',
  },
  letterSpacing: {
    50: '-0.006em',
    100: '-0.011em',
    200: '-0.014em',
    300: '-0.019em',
  },
  fontWeight: {
    normal: '400',
    medium: '600',
    strong: '700',
  },
  radii: {
    '100': '5px',
  },
};
