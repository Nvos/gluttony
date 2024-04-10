import '@fontsource/inter/400.css';
import '@fontsource/inter/500.css';
import '@fontsource/inter/700.css';
import './index.css';
import { type GlobalProvider } from '@ladle/react';

export const Provider: GlobalProvider = ({ children, globalState: { theme } }) => {
  return children;
};
