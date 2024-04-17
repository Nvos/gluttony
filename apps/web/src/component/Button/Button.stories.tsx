import type { Story } from '@ladle/react';
import { Button } from './Button';
import { css } from '~/ui/css';

export const Example: Story = () => {
  return (
    <div
      className={css({
        background: 'background.layer',
        padding: '100',
        display: 'flex',
        gap: '100',
      })}
    >
      <Button colorScheme="primary" size="md" variant="solid">
        <span>Primary/Solid</span>
      </Button>
      <Button colorScheme="primary" size="md" variant="ghost">
        Primary/Ghost
      </Button>

      <Button colorScheme="neutral" size="md" variant="solid">
        Neutral/Solid
      </Button>
      <Button colorScheme="neutral" size="md" variant="ghost">
        Neutral/Ghost
      </Button>

      <Button colorScheme="danger" size="md" variant="solid">
        Danger/Solid
      </Button>
      <Button colorScheme="danger" size="md" variant="ghost">
        Danger/Ghost
      </Button>
    </div>
  );
};
