import { ParentComponent } from 'solid-js';
import { ThemeSelector } from './ThemeSelector';

export const Layout: ParentComponent = (props) => {
  return (
    <div>
      <div>
        Side panel
        <ThemeSelector />
      </div>
      <main>{props.children}</main>
    </div>
  );
};
