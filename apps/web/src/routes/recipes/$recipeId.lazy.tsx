import { createLazyFileRoute } from '@tanstack/react-router';

export const Route = createLazyFileRoute('/recipes/$recipeId')({
  component: PostComponent,
});

function PostComponent() {
  const { recipeId } = Route.useParams();
  return <div>Recipe ID: {recipeId}</div>;
}
