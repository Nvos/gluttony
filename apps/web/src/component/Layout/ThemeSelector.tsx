import { Button } from '@gluttony/ui';
import { useEffect, useState } from 'react';
import { Sun, Moon } from 'lucide-react';

export const ThemeSelector = () => {
  const [selectedTheme, setSelectedTheme] = useState<string>('dark');
  useEffect(() => {
    document.documentElement.setAttribute('data-theme', selectedTheme);
  }, [selectedTheme]);

  const handleToggleTheme = () => {
    setSelectedTheme((prev) => {
      if (prev === 'dark') return 'light';
      return 'dark';
    });
  };

  return (
    <div>
      <Button onClick={handleToggleTheme} colorScheme="neutral" variant="ghost" size="md">
        {selectedTheme === 'light' ? <Sun /> : <Moon />}
      </Button>
    </div>
  );
};
