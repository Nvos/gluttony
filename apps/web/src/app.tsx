// @refresh reload
import { createUserClient } from '@gluttony/api';
import { darkTheme } from '@gluttony/design-system';
import { MetaProvider } from '@solidjs/meta';
import {
  cache,
  redirect,
  Route,
  Router,
  RouteLoadFunc,
  useNavigate,
  Navigate,
} from '@solidjs/router';
import { createEffect, createResource, createSignal, lazy, Suspense } from 'solid-js';
import { Layout } from './component';
import Home from './routes';
import Login from './routes/login';

const Recipes = lazy(() => import('./routes/recipes'));
const Editor = lazy(() => import('./routes/editor'));
const Recipe = lazy(() => import('./routes/recipe'));

const getSession: RouteLoadFunc = async ({}) => {
  const navigate = useNavigate();
  const api = createUserClient();
  try {
    const session = await api.me({});

    return session;
  } catch (_) {
    navigate('/login');
  }
};

export default function App() {
  // TODO: improve initial theme resolution to use user preferences + local storage
  const [selectedTheme, setSelectedTheme] = createSignal<string>(darkTheme);
  createEffect(() => {
    document.documentElement.className = selectedTheme();
    document.documentElement.dataset.theme = selectedTheme();
  });

  return (
    <Router
      root={(props) => (
        <MetaProvider>
          <Suspense>{props.children}</Suspense>
        </MetaProvider>
      )}
    >
      <Route path="/login" component={Login} />
      <Route path="/app" component={Layout} load={getSession}>
        <Route path="/" component={() => <Navigate href="recipes" />} />
        <Route path="/editor" component={Editor} />
        <Route path="/recipes" component={Recipes} />
        <Route path="/recipe/:id" component={Recipe} />
      </Route>
    </Router>
  );
}
