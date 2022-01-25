import { useLayoutEffect, useState } from 'react';

function useWindowWidth() {
  const [size, setWidth] = useState(0);
  useLayoutEffect(() => {
    function updateWidth() {
      setWidth(window.innerWidth);
    }

    window.addEventListener('resize', updateWidth);
    updateWidth();

    return () => window.removeEventListener('resize', updateWidth);
  }, []);
  return size;
}

export function DetectSmallScreenWidth() {
  return (useWindowWidth() <= 767)
}
