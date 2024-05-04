import { atom } from '../../theme';
import { Input, InputProps } from '../Input';
import { useFormField } from './useFormField';
import { forwardRef } from 'react';

export const FormInput = forwardRef<HTMLInputElement, InputProps>(
  ({ ...props }, ref) => {
    const { error, formItemId, formDescriptionId, formMessageId } =
      useFormField();

    return (
      <Input
        ref={ref}
        id={formItemId}
        aria-describedby={
          !error
            ? `${formDescriptionId}`
            : `${formDescriptionId} ${formMessageId}`
        }
        aria-invalid={!!error}
        className={atom({ marginTop: 50 })}
        {...props}
      />
    );
  },
);

FormInput.displayName = 'FormInput';
