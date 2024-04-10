import { css, cx } from '@gluttony/theme/css';
import { forwardRef, HTMLAttributes, ImgHTMLAttributes } from 'react';

const Card = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div
      ref={ref}
      className={cx(
        css({
          boxShadow: '100',
          borderRadius: '100',
          border: 'solid 1px {colors.neutral.600}',
          backgroundColor: 'neutral.100',
        }),
        className,
      )}
      {...props}
    />
  ),
);
Card.displayName = 'Card';

const CardHeader = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div
      ref={ref}
      className={cx(
        css({
          display: 'flex',
          flexDirection: 'column',
          padding: '400',
        }),
        className,
      )}
      {...props}
    />
  ),
);
CardHeader.displayName = 'CardHeader';

const CardTitle = forwardRef<HTMLParagraphElement, HTMLAttributes<HTMLHeadingElement>>(
  ({ className, ...props }, ref) => (
    <h3
      ref={ref}
      className={cx(
        css({
          textStyle: '2xl',
          fontWeight: 'medium',
        }),
        className,
      )}
      {...props}
    />
  ),
);
CardTitle.displayName = 'CardTitle';

const CardDescription = forwardRef<HTMLParagraphElement, HTMLAttributes<HTMLParagraphElement>>(
  ({ className, ...props }, ref) => (
    <p
      ref={ref}
      className={cx(
        css({
          textStyle: 'sm',
          color: 'neutral.900',
        }),
        className,
      )}
      {...props}
    />
  ),
);
CardDescription.displayName = 'CardDescription';

const CardContent = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div
      ref={ref}
      className={cx(
        css({
          padding: '400',
          paddingTop: '0',
        }),
        className,
      )}
      {...props}
    />
  ),
);
CardContent.displayName = 'CardContent';

const CardImage = forwardRef<HTMLImageElement, ImgHTMLAttributes<HTMLImageElement>>(
  ({ className, ...props }, ref) => (
    <img
      ref={ref}
      className={cx(
        css({
          display: 'block',
          backgroundSize: 'cover',
          backgroundRepeat: 'no-repeat',
          backgroundPosition: 'center',
          width: 'full',
          objectFit: 'cover',
        }),
        className,
      )}
      {...props}
    />
  ),
);
CardImage.displayName = 'CardImage';

const CardFooter = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div
      ref={ref}
      className={cx(
        css({ display: 'flex', justifyContent: 'space-between', padding: '400', paddingTop: '0' }),
        className,
      )}
      {...props}
    />
  ),
);
CardFooter.displayName = 'CardFooter';

export { Card, CardHeader, CardFooter, CardTitle, CardDescription, CardContent, CardImage };
