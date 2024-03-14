import { A, AnchorProps } from '@solidjs/router';
import { ParentComponent } from 'solid-js';
import { link, root } from './Sidebar.css';

export const Section: ParentComponent = (props) => {
  return <div></div>;
};

export const Sidebar: ParentComponent = (props) => {
  return <div class={root}>{props.children}</div>;
};

export const Link: ParentComponent<AnchorProps> = (props) => {
  return (
    <A
      {...props}
      end
      classList={{
        [link]: true,
      }}
    >
      {props.children}
    </A>
  );
};
