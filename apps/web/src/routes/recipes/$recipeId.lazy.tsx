import { createLazyFileRoute } from '@tanstack/react-router';
import { RecipeView } from '~/view/RecipeView/RecipeView';

const View = () => {
  const { recipeId } = Route.useParams();

  return <RecipeView recipeId={recipeId} />;
};

export const Route = createLazyFileRoute('/recipes/$recipeId')({
  component: View,
});
