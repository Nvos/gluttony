import { RecipeEditor } from '~/component/RecipeEditor';

type RecipeViewProps = {
  recipeId: number;
};

export const RecipeView = ({ recipeId }: RecipeViewProps) => {
  return (
    <div>
      Recipe editor {recipeId}
      <RecipeEditor />
    </div>
  );
};
