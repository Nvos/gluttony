import { createEffect, createSignal } from 'solid-js';
import { lightTheme } from '@gluttony/design-system';

export const ThemeSelector = () => {
  // Resolve initial theme depending on system and local storage
  const [selectedTheme, setSelectedTheme] = createSignal<string>(lightTheme);
  createEffect(() => {
    document.documentElement.className = selectedTheme();
    document.documentElement.dataset.theme = selectedTheme();
  });

  // Later on toggle component
  return null;
};
