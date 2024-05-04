import '@gluttony/ui/global.css';
import { Provider as JotaiProvider } from 'jotai';
import { routeTree } from './routeTree.gen';
import { RouterProvider, createRouter } from '@tanstack/react-router';
import { ThemeProvider } from '@gluttony/ui/ThemeProvider';

const router = createRouter({
  routeTree,
  context: {
    auth: undefined!,
  },
});

// Register the router instance for type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}

export default function App() {
  return (
    <ThemeProvider>
      <JotaiProvider>
        <RouterProvider router={router} />
      </JotaiProvider>
    </ThemeProvider>
  );
}
