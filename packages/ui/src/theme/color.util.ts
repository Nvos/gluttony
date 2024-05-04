export const radixToLocalTokenScale = (value: { [key: string]: string }) => {
  const name = Object.keys(value)
    .find((it) => it.endsWith('12'))!
    .replace('12', '');

  const getByScale = (scale: number) => {
    return value[`${name}${scale}`]!;
  };

  return {
    50: getByScale(1),
    100: getByScale(2),
    200: getByScale(3),
    300: getByScale(4),
    400: getByScale(5),
    500: getByScale(6),
    600: getByScale(7),
    700: getByScale(8),
    800: getByScale(9),
    900: getByScale(10),
    950: getByScale(11),
    1000: getByScale(12),
  };
};
