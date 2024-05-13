import { Button } from '../Button';
import { Popover, PopoverContent, PopoverTrigger } from './Popover';

export const Example = () => {
  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="primary">Open popover</Button>
      </PopoverTrigger>
      <PopoverContent>
        <div>Content</div>
      </PopoverContent>
    </Popover>
  );
};
