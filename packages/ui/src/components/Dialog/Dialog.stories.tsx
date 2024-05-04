import { useForm } from 'react-hook-form';
import { atom } from '../../theme';
import { Button } from '../Button';
import {
  Form,
  FormDescription,
  FormInput,
  FormItem,
  FormLabel,
  FormMessage,
} from '../Form';
import { Input } from '../Input';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from './Dialog';

type Schema = {
  username: string;
};

export const Example = () => {
  const form = useForm<Schema>({
    defaultValues: {
      username: '',
    },
  });

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="primary">Open</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Edit profile</DialogTitle>
          <DialogDescription>
            Make changes to your profile here. Click save when you're done.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form
            className={atom({
              paddingY: 300,
              display: 'flex',
              flexDirection: 'column',
            })}
          >
            <FormItem>
              <FormLabel>Username</FormLabel>
              <FormInput />
              <FormDescription>
                This is your public display name.
              </FormDescription>
            </FormItem>
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormInput />
            </FormItem>
          </form>
        </Form>
        <DialogFooter>
          <Button type="submit" variant="primary">
            Save changes
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
