import { createFileRoute } from '@tanstack/react-router';
import { z } from 'zod';

export const Route = createFileRoute('/recipes/$recipeId')({
  
  parseParams: (params) => ({
    recipeId: z.number().int().parse(Number(params.recipeId)),
  }),
});
