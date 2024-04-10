import { css } from '@gluttony/theme/css';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardImage,
  CardTitle,
} from '@gluttony/ui';
import BigosImg from '~/assets/bigos.jpg';
import { Link } from '@tanstack/react-router';

const cardLink = css({
  width: '[320px]',
});

export const RecipeView = () => {
  return (
    <div
      className={css({
        padding: '500',
        minHeight: 'dvh',
        width: 'full',
      })}
    >
      <Link
        className={css({
          display: 'inline-flex',
        })}
        to="/recipes/$recipeId"
        params={{
          recipeId: 1,
        }}
      >
        <Card className={cardLink}>
          <CardHeader>
            <CardTitle>Old polish hunter&apos;s stew</CardTitle>
          </CardHeader>
          <CardImage src={BigosImg} />
          <CardContent />
          <CardFooter>Footer</CardFooter>
        </Card>
      </Link>
    </div>
  );
};
