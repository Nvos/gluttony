import { Input, InputProps } from '../Input';
import { useFormField } from './Form';
import { forwardRef } from 'react';
import { css } from '~/ui/css';

export const FormInput = forwardRef<HTMLInputElement, InputProps>(({ ...props }, ref) => {
  const { error, formItemId, formDescriptionId, formMessageId } = useFormField();

  return (
    <Input
      ref={ref}
      id={formItemId}
      aria-describedby={!error ? `${formDescriptionId}` : `${formDescriptionId} ${formMessageId}`}
      aria-invalid={!!error}
      className={css({ marginTop: '50' })}
      {...props}
    />
  );
});

FormInput.displayName = 'FormInput';