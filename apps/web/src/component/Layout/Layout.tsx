import { PropsWithChildren } from 'react';
import { ThemeSelector } from './ThemeSelector';
import { Link } from '@tanstack/react-router';
import { Beef, CookingPot } from 'lucide-react';
import { Button } from '@gluttony/ui/Button';
import { atom, vars } from '@gluttony/ui';
import { linkRoot } from './Link.css';

export const Layout = ({ children }: PropsWithChildren<unknown>) => {
  return (
    <div
      className={atom({
        display: 'flex',
      })}
    >
      <aside
        className={atom({
          paddingX: 400,
          paddingY: 500,
          display: 'flex',
          flexDirection: 'column',
          gap: 600,
          minHeight: 'dvh',
        })}
        style={{
          width: 240,
          backgroundColor: vars.color.surface[100],
          borderRight: `solid 1px ${vars.color.neutral[500]}`,
        }}
      >
        <nav
          className={atom({
            display: 'flex',
            flexDirection: 'column',
            gap: 100,
          })}
        >
          <li>
            <Link className={linkRoot} to="/">
              <Beef /> Gluttony
            </Link>
          </li>
          <li>
            <Link className={linkRoot} to="/recipes">
              <CookingPot />
              Recipes
            </Link>
          </li>
        </nav>
        <div className={atom({ flex: 1 })}></div>
        <div>
          <ThemeSelector />
          <Button variant="secondary">Logout</Button>
        </div>
      </aside>
      <main className={atom({ width: 'full' })}>{children}</main>
    </div>
  );
};
