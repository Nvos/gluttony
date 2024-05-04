import { ReactNode, forwardRef, ButtonHTMLAttributes } from 'react';
import { root, ButtonVariants } from './Button.css';
import { clsx } from 'clsx';

export type ButtonProps = {
  children?: ReactNode;
} & ButtonHTMLAttributes<HTMLButtonElement> &
  ButtonVariants;

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ children, className, size, variant, type = 'button', ...rest }, ref) => {
    return (
      <button
        ref={ref}
        className={clsx(root({ size: size, variant: variant }))}
        type={type}
        {...rest}
      >
        {children}
      </button>
    );
  },
);
