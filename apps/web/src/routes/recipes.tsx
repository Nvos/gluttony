import { Title } from '@solidjs/meta';
import { BadRequest, createRecipeClient } from '@gluttony/api';
import { createEffect, createMemo, createResource, For, Suspense } from 'solid-js';
import { createForm, SubmitHandler, insert, getValues, FormError } from '@modular-forms/solid';
import { createSignal } from 'solid-js';
import { createScheduled, debounce } from '@solid-primitives/scheduled';
import { Button } from '@gluttony/ui';
import { A } from '@solidjs/router';
import { ConnectError } from '@connectrpc/connect';

type RecipeStep = {
  order: number;
  description: string;
};

type Recipe = {
  name: string;
  description: string;
  steps: RecipeStep[];
};

export default function Home() {
  const [search, setSearch] = createSignal('');
  const scheduled = createScheduled((fn) => debounce(fn, 1000));
  const [recipeForm, { Form, Field, FieldArray }] = createForm<Recipe>({
    initialValues: {
      description: '',
      name: '',
      steps: [],
    },
  });
  const api = createRecipeClient();

  const debouncedSearch = createMemo((prev: string = '') => {
    const value = search();
    return scheduled() ? value : prev;
  });

  createEffect(() => {
    console.log('search', debouncedSearch());
  });
  const [recipes, { refetch }] = createResource(debouncedSearch, (props) =>
    api.allRecipes({ limit: 20, offset: 0, search: props }),
  );

  const handleSubmit: SubmitHandler<Recipe> = (data) => {
    console.log(data);

    if (data.name === 'admin') {
      throw new FormError<Recipe>('Validation error.', {
        name: 'not unique',
      });
    }

    api
      .createRecipe(data)
      .then((it) => {
        console.log('got recipeId', it.recipeId);
        refetch();
      })
      .catch((err) => {
        if (err instanceof ConnectError) {
          const connectErr = ConnectError.from(err);
          console.log('connectErr: details', connectErr.details);
          const badRequest = err.findDetails(BadRequest);
          console.log('connectErr: rawMessage', connectErr.rawMessage);
          console.log('connectErr: badRequest', badRequest);
        }
      });
  };

  return (
    <main>
      <Title>Recipes</Title>
      <h1>Recipes</h1>
      <label>
        Search by name
        <input onInput={(event) => setSearch(event.currentTarget.value)} />
      </label>
      <Form
        onSubmit={handleSubmit}
        style={{
          display: 'flex',
          'flex-direction': 'column',
          gap: '16px',
          'max-width': '560px',
        }}
      >
        <Field name="name">
          {(field, props) => (
            <label>
              Recipe name
              <input {...props} />
              {field.error && <div>{field.error}</div>}
            </label>
          )}
        </Field>
        <Field name="description">
          {(_, props) => (
            <label>
              Recipe description
              <input {...props} />
            </label>
          )}
        </Field>
        <FieldArray name="steps">
          {(fields) => (
            <For each={fields.items}>
              {(_, index) => (
                <div
                  style={{
                    display: 'flex',
                    'flex-direction': 'column',
                    gap: '8px',
                  }}
                >
                  <Field name={`steps.${index()}.order`} type="number">
                    {(field, props) => (
                      <label>
                        Step order
                        <input {...props} type="number" />
                      </label>
                    )}
                  </Field>
                  <Field name={`steps.${index()}.description`}>
                    {(field, props) => (
                      <label>
                        Step description
                        <input {...props} />
                      </label>
                    )}
                  </Field>
                </div>
              )}
            </For>
          )}
        </FieldArray>
        <Button
          variant="outline"
          colorScheme="background"
          type="button"
          onClick={() => {
            const value = getValues(recipeForm);
            let order = 1;
            if (value.steps !== undefined && value.steps.length > 0) {
              const latestOrder = value.steps[value.steps.length - 1]?.order ?? 0;
              order = latestOrder + 1;
            }

            insert(recipeForm, 'steps', { value: { order: order, description: '' } });
          }}
        >
          Add step
        </Button>
        <Button variant="outline" colorScheme="primary" type="submit">
          Save
        </Button>
      </Form>

      <Suspense fallback={'Loading recipes...'}>
        <For each={recipes()?.recipes ?? []}>
          {(value) => (
            <div>
              <div>{value.id} </div>
              <div>{value.name} </div>
              <div>{value.description} </div>
              <A href={`/recipe/${value.id}`}>Goto</A>
            </div>
          )}
        </For>
      </Suspense>
    </main>
  );
}
