import type { Story } from '@ladle/react';
import { css } from '~/ui/css';
import { Input } from './Input';

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
      <Input />
    </div>
  );
};
