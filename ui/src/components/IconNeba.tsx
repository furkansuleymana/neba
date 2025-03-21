interface IconNebaProps extends React.ComponentPropsWithoutRef<"svg"> {
  size?: number | string;
}

export function IconNeba({ size }: IconNebaProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 100 100"
    >
      <rect width="100" height="100" rx="20" fill="#eeeeee" />
      <path
        d="M47.17 74.21L23.48 74.21L23.48 67.45L30.07 66.13L30.07 34.79L23.13 33.48L23.13 26.67L40.27 26.67L40.84 33.56Q43.21 29.87 46.75 27.83Q50.29 25.79 54.77 25.79L54.77 25.79Q62.11 25.79 66.19 30.29Q70.28 34.79 70.28 44.38L70.28 44.38L70.28 66.13L76.87 67.45L76.87 74.21L53.01 74.21L53.01 67.45L59.21 66.13L59.21 44.46Q59.21 39.06 57.05 36.79Q54.90 34.53 50.64 34.53L50.64 34.53Q47.43 34.53 45.03 35.96Q42.64 37.39 41.15 39.94L41.15 39.94L41.15 66.13L47.17 67.45L47.17 74.21Z"
        fill="#222831"
      />
    </svg>
  );
}
