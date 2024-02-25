import { ComponentProps } from 'solid-js';
import { root } from './Button.css';

export const Button = (props: ComponentProps<'button'>) => {
  return <button class={root}>{props.children}</button>;
};
