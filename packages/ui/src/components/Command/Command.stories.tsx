import { atom } from '../../theme';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
  CommandShortcut,
} from './Command';
import {
  Calculator,
  Calendar,
  CreditCard,
  Settings,
  Smile,
  User,
} from 'lucide-react';

export const Example = () => {
  return (
    <Command
      className={atom({ borderRadius: 100, boxShadow: 100, border: 'neutral' })}
    >
      <CommandInput placeholder="Type a command or search..." />
      <CommandList>
        <CommandEmpty>No results found.</CommandEmpty>
        <CommandGroup heading="Suggestions">
          <CommandItem>
            <Calendar
              className={atom({ marginRight: 100, width: 25, height: 25 })}
            />
            <span>Calendar</span>
          </CommandItem>
          <CommandItem>
            <Smile
              className={atom({ marginRight: 100, width: 25, height: 25 })}
            />
            <span>Search Emoji</span>
          </CommandItem>
          <CommandItem>
            <Calculator
              className={atom({ marginRight: 100, width: 25, height: 25 })}
            />
            <span>Calculator</span>
          </CommandItem>
        </CommandGroup>
        <CommandSeparator />
        <CommandGroup heading="Settings">
          <CommandItem>
            <User
              className={atom({ marginRight: 100, width: 25, height: 25 })}
            />
            <span>Profile</span>
            <CommandShortcut>⌘P</CommandShortcut>
          </CommandItem>
          <CommandItem>
            <CreditCard
              className={atom({ marginRight: 100, width: 25, height: 25 })}
            />
            <span>Billing</span>
            <CommandShortcut>⌘B</CommandShortcut>
          </CommandItem>
          <CommandItem>
            <Settings
              className={atom({ marginRight: 100, width: 25, height: 25 })}
            />
            <span>Settings</span>
            <CommandShortcut>⌘S</CommandShortcut>
          </CommandItem>
        </CommandGroup>
      </CommandList>
    </Command>
  );
};
