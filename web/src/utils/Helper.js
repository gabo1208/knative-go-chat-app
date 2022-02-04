import { useLayoutEffect, useState } from 'react';

export const NewUserConnected = "gabo1208.go-chat-client.NewUserConnected"
export const UserDisconnected = "gabo1208.go-chat-client.UserDisconnected"
export const ChatAppEventSource = "gabo1208.go-chat-client/source"

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
