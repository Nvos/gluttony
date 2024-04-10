import { ButtonHTMLAttributes, ReactNode, forwardRef } from 'react';
import { button } from '@gluttony/theme/recipes';
import { css, cx } from '@gluttony/theme/css';
import { type ColorPalette } from '@gluttony/theme/tokens';
import { type RecipeVariantProps } from '@gluttony/theme/types';

type ButtonVariants = RecipeVariantProps<typeof button>;

type Props = {
  colorScheme?: ColorPalette;
} & ButtonVariants &
  ButtonHTMLAttributes<HTMLButtonElement>;

export const Button = forwardRef<HTMLButtonElement, Props>(
  ({ colorScheme, size, variant, className, children, ...rest }, ref) => {
    return (
      <button
        ref={ref}
        className={cx(
          css({ colorPalette: colorScheme }),
          button({
            size,
            variant,
          }),
          className,
        )}
        {...rest}
      >
        {children}
      </button>
    );
  },
);
