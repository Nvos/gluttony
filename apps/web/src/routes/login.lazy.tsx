import { css } from '~/ui/theme/css';
import { createLazyFileRoute } from '@tanstack/react-router';

export const Route = createLazyFileRoute('/login')({
  component: About,
});

function About() {
  return (
    <div className={css({ display: 'flex', color: 'primary.900' })}>Hello wtf from About!</div>
  );
}
