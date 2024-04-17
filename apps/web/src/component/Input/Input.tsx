import { forwardRef, type InputHTMLAttributes } from 'react';
import { cx } from '~/ui/css';
import { input } from '~/ui/recipes';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {}

const Input = forwardRef<HTMLInputElement, InputProps>(({ className, type, ...props }, ref) => {
  return (
    <input type={type} className={cx(className, input({ size: 'md' }))} ref={ref} {...props} />
  );
});
Input.displayName = 'Input';

export { Input };
