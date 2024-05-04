import { atom } from '../../theme';
import { Button } from './Button';

export const Example = () => {
  return (
    <div className={atom({ display: 'flex', gap: 200 })}>
      <Button variant="primary">Primary</Button>
      <Button variant="secondary">Secondary</Button>
      <Button variant="destructive">Destructive</Button>
    </div>
  );
};
