import { Button } from '@gluttony/ui/Button';
import { Sun, Moon } from 'lucide-react';
import { useTheme } from '@gluttony/ui/ThemeProvider';

export const ThemeSelector = () => {
  const { theme, toggleTheme } = useTheme();

  return (
    <div>
      <Button onClick={toggleTheme} variant="secondary" size="md">
        {theme === 'light' ? <Sun /> : <Moon />}
      </Button>
    </div>
  );
};
