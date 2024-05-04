import { GlobalProvider, ThemeState } from '@ladle/react';
import '@fontsource/inter/400.css';
import '@fontsource/inter/600.css';
import '@fontsource/inter/700.css';

import '../src/theme/global.css';
import { darkTheme } from '../src/theme/theme.dark.css';
import { lightTheme } from '../src/theme/theme.light.css';
import { vars } from '../src/theme';

export const Provider: GlobalProvider = ({ children, globalState }) => (
  <div
    className={globalState.theme === ThemeState.Dark ? darkTheme : lightTheme}
    style={{
      backgroundColor: vars.color.surface[50],
      padding: vars.space[600],
    }}
  >
    <div>{children}</div>
  </div>
);
