// @refresh reload
import { MetaProvider, Title } from '@solidjs/meta';
import { Router } from '@solidjs/router';
import { FileRoutes } from '@solidjs/start';
import { Suspense } from 'solid-js';
import { Layout } from './component';

export default function App() {
  return (
    <Router
      root={(props) => (
        <MetaProvider>
          <Title>SolidStart - Basic</Title>
          <a href="/">Index</a>
          <a href="/about">About</a>
          <Layout>
            <Suspense>{props.children}</Suspense>
          </Layout>
        </MetaProvider>
      )}
    >
      <FileRoutes />
    </Router>
  );
}
