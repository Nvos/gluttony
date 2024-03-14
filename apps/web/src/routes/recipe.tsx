import { createRecipeClient } from '@gluttony/api';
import { useParams, RouteSectionProps } from '@solidjs/router';
import { Component, createEffect } from 'solid-js';

type RouteParams = {
  id: string;
};

const Recipe: Component<RouteSectionProps> = () => {
  const params = useParams<RouteParams>();
  const api = createRecipeClient();

  createEffect(() => {
    api
      .singleRecipe({ id: parseInt(params.id, 10) })
      .then(console.log)
      .catch(console.error);
  });

  return <div>Recipe wow, {params.id}</div>;
};

export default Recipe;
