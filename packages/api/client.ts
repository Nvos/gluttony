import { createPromiseClient } from '@connectrpc/connect';
import { RecipeService } from './recipe/v1/recipe_connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { ServiceType } from '@bufbuild/protobuf';
import { UserService } from './user/v1/user_connect';

const transport = createConnectTransport({
  // TODO(AK) 05/03/2024: customizable baseUrl via env
  baseUrl: 'http://localhost:6001',
  credentials: 'include',
});

export const createClient = <T extends ServiceType>(service: T) => {
  return createPromiseClient(service, transport);
};

export const createRecipeClient = () => createClient(RecipeService);
export const createUserClient = () => createClient(UserService);
