import { createLazyFileRoute } from '@tanstack/react-router';
import { RecipeEditor } from '~/component';
import { css } from '~/ui/css';

export const Route = createLazyFileRoute('/recipes/$recipeId')({
  component: PostComponent,
});

function PostComponent() {
  const { recipeId } = Route.useParams();
  return (
    <div
      className={css({
        padding: '300',
      })}
    >
      <p>Recipe ID: {recipeId}</p>
      <RecipeEditor />
    </div>
  );
}
