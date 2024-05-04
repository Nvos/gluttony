import type { Story } from '@ladle/react';
import { useForm } from 'react-hook-form';
import { FormInput } from './FormInput';
import { Button } from '../Button/Button';
import {
  FormField,
  FormItem,
  FormLabel,
  FormDescription,
  FormMessage,
  Form,
} from './Form';
import { atom } from '../../theme';

type Schema = {
  username: string;
};

export const Example: Story = () => {
  const form = useForm<Schema>({
    defaultValues: {
      username: '',
    },
  });
  return (
    <div
      className={atom({
        padding: 100,
        display: 'flex',
        gap: 100,
      })}
    >
      <Form {...form}>
        <form>
          <FormField
            control={form.control}
            name="username"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Username</FormLabel>
                <FormInput {...field} />
                <FormDescription>
                  This is your public display name.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button variant="primary" type="submit">
            Submit
          </Button>
        </form>
      </Form>
    </div>
  );
};
