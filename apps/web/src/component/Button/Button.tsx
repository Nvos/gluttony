import { ButtonHTMLAttributes, forwardRef } from 'react';
import { button } from '~/ui/recipes';
import { css, cx } from '~/ui/css';
import { type ColorPalette } from '~/ui/tokens';
import { type RecipeVariantProps } from '~/ui/types';

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
Button.displayName = 'Button';
