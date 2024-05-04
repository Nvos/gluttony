import { GlobalProvider, ThemeState } from '@ladle/react';
import '@fontsource/inter/400.css';
import '@fontsource/inter/600.css';
import '@fontsource/inter/700.css';

import '../src/theme/global.css';
import { darkTheme } from '../src/theme/theme.dark.css';
import { lightTheme } from '../src/theme/theme.light.css';
import { vars } from '../src/theme';
import { ThemeProvider, useTheme } from '../src/components/ThemeProvider';
import { useEffect } from 'react';

export const Provider: GlobalProvider = ({ children, globalState }) => (
  <ThemeProvider>
    <ThemeSync theme={globalState.theme} />
    <div
      style={{
        backgroundColor: vars.color.surface[50],
        padding: vars.space[600],
      }}
    >
      <div>{children}</div>
    </div>
  </ThemeProvider>
);

const ThemeSync = ({ theme }: { theme: ThemeState }) => {
  const { setTheme } = useTheme();
  useEffect(() => {
    setTheme(theme === ThemeState.Dark ? 'dark' : 'light');
  }, [theme]);

  return null;
};
