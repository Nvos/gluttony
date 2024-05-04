import { Card, CardContent, CardFooter, CardHeader, CardImage, CardTitle } from '@gluttony/ui/Card';
import { atom } from '@gluttony/ui';
import BigosImg from '~/assets/bigos.jpg';
import { Link } from '@tanstack/react-router';
import { card } from './RecipeView.css';

export const RecipeView = () => {
  return (
    <div
      className={atom({
        padding: 500,
        minHeight: 'dvh',
        width: 'full',
      })}
    >
      <Link
        className={atom({
          display: 'inline-flex',
        })}
        to="/recipes/$recipeId"
        params={{
          recipeId: 1,
        }}
      >
        <Card className={card}>
          <CardHeader>
            <CardTitle>Old polish hunter&apos;s stew</CardTitle>
          </CardHeader>
          <div className={atom({ position: 'relative', paddingBottom: 400 })}>
            <CardImage src={BigosImg} />
          </div>
          <CardContent>
            <div>Meat</div>
          </CardContent>
          <CardFooter></CardFooter>
        </Card>
      </Link>
    </div>
  );
};
