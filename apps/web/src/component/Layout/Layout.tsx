import { createUserClient } from '@gluttony/api';
import { Button } from '@gluttony/ui';
import { useNavigate } from '@solidjs/router';
import { ParentComponent } from 'solid-js';
import { root } from './Layout.css';
import { Link, Sidebar } from './Sidebar';
import { ThemeSelector } from './ThemeSelector';

export const Layout: ParentComponent = (props) => {
  const userApi = createUserClient();
  const navigation = useNavigate();
  const handleLogout = () =>
    userApi.logout({}).then(() => {
      navigation('/login');
    });

  return (
    <div class={root}>
      <Sidebar>
        <ThemeSelector />

        <Link href="/app/recipes">Recipes</Link>
        <Button onClick={handleLogout}>Logout</Button>
      </Sidebar>
      <main>{props.children}</main>
    </div>
  );
};
