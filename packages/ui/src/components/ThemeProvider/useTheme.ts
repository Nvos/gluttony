import { createContext, useContext } from 'react';

export type Theme = 'dark' | 'light';
export type ThemeContextState = {
  theme: Theme;
  toggleTheme: () => void;
  setTheme: (theme: Theme) => void;
};

export const ThemeContext = createContext({} as ThemeContextState);

export const useTheme = () => {
  const ctx = useContext(ThemeContext);
  if (ctx === undefined) {
    throw new Error(
      'Theme context value is undefined, context provider might be missing',
    );
  }

  return ctx;
};
