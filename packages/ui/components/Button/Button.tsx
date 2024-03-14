import { ComponentProps } from 'solid-js';
import { RecipeVariants } from '@vanilla-extract/recipes';
import { root } from './Button.css';
import { splitProps } from 'solid-js';

type Variants = RecipeVariants<typeof root>;

export const Button = (props: ComponentProps<'button'> & Variants) => {
  const [local, rest] = splitProps(props, ['variant', 'colorScheme']);

  return (
    <button
      type="button"
      class={root({ variant: local.variant, colorScheme: local.colorScheme })}
      {...rest}
    >
      {props.children}
    </button>
  );
};
