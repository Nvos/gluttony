import { PropsWithChildren } from 'react';
import { ThemeSelector } from './ThemeSelector';
import { css } from '@gluttony/theme/css';
import { Link } from '@tanstack/react-router';
import { Beef, CookingPot } from 'lucide-react';
import { Button } from '@gluttony/ui';

const linkStyles = css({
  display: 'flex',
  gap: '200',
  alignItems: 'center',
  paddingX: '200',
  height: '[40px]',
  width: 'full',
  borderRadius: '100',
  transition: 'colors',
  '&.active': {
    color: 'primary.950',
    fontWeight: 'heavy',
  },
  '&:hover': {
    backgroundColor: 'neutral.300',
  },
});

const liStyles = css({
  listStyle: 'none',
});

const navStyles = css({
  display: 'flex',
  gap: '100',
  flexDirection: 'column',
});

export const Layout = ({ children }: PropsWithChildren<unknown>) => {
  return (
    <div
      className={css({
        display: 'flex',
      })}
    >
      <aside
        className={css({
          paddingX: '400',
          paddingY: '500',
          display: 'flex',
          flexDirection: 'column',
          gap: '600',
          minHeight: 'dvh',
          width: '[240px]',
          backgroundColor: 'neutral.100',
          borderRight: '{colors.neutral.500} solid 1px',
        })}
      >
        <nav className={navStyles}>
          <li className={liStyles}>
            <Link className={linkStyles} to="/">
              <Beef /> Gluttony
            </Link>
          </li>
          <li className={liStyles}>
            <Link className={linkStyles} to="/recipes">
              <CookingPot />
              Recipes
            </Link>
          </li>
        </nav>
        <div className={css({ flex: '1' })}></div>
        <div>
          <ThemeSelector />
          <Button variant="ghost" colorScheme="neutral">
            Logout
          </Button>
        </div>
      </aside>
      <main className={css({ width: 'full' })}>{children}</main>
    </div>
  );
};
