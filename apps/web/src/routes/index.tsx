import { Button } from '@gluttony/ui';
import { Title } from '@solidjs/meta';

export default function Home() {
  return (
    <main>
      <Title>Hello World</Title>
      <h1>Hello world!</h1>
      <p>
        Visit{' '}
        <a href="https://start.solidjs.com" target="_blank">
          start.solidjs.com
        </a>{' '}
        to learn how to build SolidStart apps.
      </p>
      <Button>We have button there</Button>
    </main>
  );
}
