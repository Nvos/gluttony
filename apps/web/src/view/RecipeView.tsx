import { css } from '~/ui/css';
import {
  Button,
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardImage,
  CardTitle,
} from '~/component';
import BigosImg from '~/assets/bigos.jpg';
import { Link } from '@tanstack/react-router';
import { CookingPot } from 'lucide-react';

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
          <div className={css({ position: 'relative', paddingBottom: '400' })}>
            <CardImage src={BigosImg} />
            <Button
              size="icon"
              colorScheme="primary"
              variant="solid"
              className={css({
                borderRadius: 'round',
                position: 'absolute',
                right: '300',
                bottom: '500',
              })}
            >
              <CookingPot />
            </Button>
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
