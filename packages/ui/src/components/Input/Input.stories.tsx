import { atom } from '../../theme';
import { Input } from './Input';

export const Example = () => {
  return (
    <div className={atom({ display: 'flex', gap: 200 })}>
      <Input />
    </div>
  );
};
