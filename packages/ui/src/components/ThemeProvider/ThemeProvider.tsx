import { ReactNode, useEffect, useState } from 'react';
import { darkTheme, lightTheme } from '../../theme';
import { type Theme, ThemeContext } from './useTheme';

const STORAGE_KEY = 'gluttony-theme';

type ThemeProviderProps = {
  children?: ReactNode;
};

const resolveTheme = () => {
  const current = localStorage.getItem(STORAGE_KEY);
  if (current !== null) return current as Theme;

  const isDark = matchMedia('(prefers-color-scheme: dark)').matches;

  return isDark ? 'dark' : 'light';
};

export const ThemeProvider = ({ children }: ThemeProviderProps) => {
  const [theme, setTheme] = useState<Theme>(resolveTheme);
  const themeClass = theme === 'dark' ? darkTheme : lightTheme;

  useEffect(() => {
    document.documentElement.classList.remove(lightTheme, darkTheme);
    document.documentElement.classList.add(themeClass);
  }, [themeClass]);

  const handleThemeToggle = () => {
    let nextTheme = theme === 'dark' ? ('light' as const) : ('dark' as const);
    localStorage.setItem(STORAGE_KEY, nextTheme);
    setTheme(nextTheme);
  };

  const handleSetTheme = (theme: Theme) => {
    localStorage.setItem(STORAGE_KEY, theme);
    setTheme(theme);
  };

  return (
    <ThemeContext.Provider
      value={{
        theme: theme,
        toggleTheme: handleThemeToggle,
        setTheme: handleSetTheme,
      }}
    >
      {children}
    </ThemeContext.Provider>
  );
};
