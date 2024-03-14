import { createEffect, createSignal } from 'solid-js';
import { lightTheme, darkTheme } from '@gluttony/design-system';
import { Button } from '@gluttony/ui';

export const ThemeSelector = () => {
  // Resolve initial theme depending on system and local storage
  const [selectedTheme, setSelectedTheme] = createSignal<string>(darkTheme);
  createEffect(() => {
    document.documentElement.className = selectedTheme();
    document.documentElement.dataset.theme = selectedTheme();
  });

  const handleToggleTheme = () => {
    if (selectedTheme() == darkTheme) {
      setSelectedTheme(lightTheme);
      return;
    }

    setSelectedTheme(darkTheme);
  };

  // Later on toggle component
  return <Button onClick={handleToggleTheme}>Toggle</Button>;
};
