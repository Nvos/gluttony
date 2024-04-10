import { useEffect, useState } from 'react';
import { Provider as JotaiProvider } from 'jotai';

// const getSession: RouteLoadFunc = async ({}) => {
//   const navigate = useNavigate();
//   const api = createUserClient();
//   try {
//     const session = await api.me({});

//     return session;
//   } catch (_) {
//     navigate('/login');
//   }
// };

import { routeTree } from './routeTree.gen';
import { RouterProvider, createRouter } from '@tanstack/react-router';
// import { authAtom } from './state/auth';

// Create a new router instance
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

// const RootApp = () => {
//   const auth = useAtomValue(authAtom);
// }

export default function App() {
  // TODO: improve initial theme resolution to use user preferences + local storage
  const [selectedTheme] = useState<string>('');
  useEffect(() => {
    document.documentElement.className = selectedTheme;
    document.documentElement.dataset.theme = selectedTheme;
  }, [selectedTheme]);

  return (
    <JotaiProvider>
      <RouterProvider router={router} />
    </JotaiProvider>
    // <Router
    //   root={(props) => (
    //     <MetaProvider>
    //       <Suspense>{props.children}</Suspense>
    //     </MetaProvider>
    //   )}
    // >
    //   <Route path="/login" component={Login} />
    //   <Route path="/app" component={Layout} load={getSession}>
    //     <Route path="/" component={() => <Navigate href="recipes" />} />
    //     <Route path="/editor" component={Editor} />
    //     <Route path="/recipes" component={Recipes} />
    //     <Route path="/recipe/:id" component={Recipe} />
    //   </Route>
    // </Router>
  );
}
