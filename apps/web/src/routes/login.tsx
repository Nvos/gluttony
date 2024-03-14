import { createUserClient } from '@gluttony/api';
import { Button } from '@gluttony/ui';
import { createForm, SubmitHandler } from '@modular-forms/solid';
import { useNavigate } from '@solidjs/router';
import { Show, createSignal } from 'solid-js';

type LoginInput = {
  username: string;
  password: string;
};

const Login = () => {
  const [isInvalid, setInvalid] = createSignal(false);
  const [_, { Form, Field }] = createForm<LoginInput>({
    initialValues: {
      password: '',
      username: '',
    },
  });

  const api = createUserClient();
  const navigate = useNavigate();

  const handleSubmit: SubmitHandler<LoginInput> = (data) => {
    api
      .login({ username: data.username, password: data.password })
      .then(() => {
        navigate('/app');
      })
      .catch(() => setInvalid(true));
  };

  return (
    <Form onSubmit={handleSubmit}>
      <label style="display: block">
        Username
        <Field name="username">{(_, props) => <input {...props} />}</Field>
      </label>

      <label style="display: block">
        Password
        <Field name="password">{(_, props) => <input {...props} />}</Field>
      </label>

      <Show when={isInvalid()}>Invalid credentials</Show>
      <Button type="submit" colorScheme="primary" variant="solid">
        Login
      </Button>
    </Form>
  );
};

export default Login;
