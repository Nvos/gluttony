import { createLazyFileRoute } from '@tanstack/react-router';
import { RecipeView } from '~/view/RecipeView';

export const Route = createLazyFileRoute('/recipes/')({
  component: RecipeView,
});
