export const radixToLocalTokenScale = (
  light: { [key: string]: string },
  dark: { [key: string]: string },
) => {
  const name = Object.keys(light)
    .find((it) => it.endsWith('12'))!
    .replace('12', '');

  const getByScale = (color: { [key: string]: string }, scale: number) => {
    return {
      value: color[`${name}${scale}`]!,
    };
  };

  return {
    light: {
      50: getByScale(light, 1),
      100: getByScale(light, 2),
      200: getByScale(light, 3),
      300: getByScale(light, 4),
      400: getByScale(light, 5),
      500: getByScale(light, 6),
      600: getByScale(light, 7),
      700: getByScale(light, 8),
      800: getByScale(light, 9),
      900: getByScale(light, 10),
      950: getByScale(light, 11),
      1000: getByScale(light, 12),
    },
    dark: {
      50: getByScale(dark, 1),
      100: getByScale(dark, 2),
      200: getByScale(dark, 3),
      300: getByScale(dark, 4),
      400: getByScale(dark, 5),
      500: getByScale(dark, 6),
      600: getByScale(dark, 7),
      700: getByScale(dark, 8),
      800: getByScale(dark, 9),
      900: getByScale(dark, 10),
      950: getByScale(dark, 11),
      1000: getByScale(dark, 12),
    },
  };
};

type LocalColor = ReturnType<typeof radixToLocalTokenScale>;

export const colorTokenToSemanticToken = (name: string) => {
  const getValue = (scale: keyof LocalColor['light']) => {
    return {
      value: { base: `{colors.${name}.light.${scale}}`, _dark: `{colors.${name}.dark.${scale}}` },
    };
  };

  return {
    50: getValue(50),
    100: getValue(100),
    200: getValue(200),
    300: getValue(300),
    400: getValue(400),
    500: getValue(500),
    600: getValue(600),
    700: getValue(700),
    800: getValue(800),
    900: getValue(900),
    950: getValue(950),
    1000: getValue(1000),
  };
};
