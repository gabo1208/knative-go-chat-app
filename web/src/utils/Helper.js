import { useLayoutEffect, useState } from 'react';

export const NewUserConnected = "gabo1208.knative-go-chat-app.NewUserConnected"
export const FirstUserConnection = "gabo1208.knative-go-chat-app.FirstUserConnection"
export const UserReconnected = "gabo1208.knative-go-chat-app.UserReconnected"
export const UserDisconnected = "gabo1208.knative-go-chat-app.UserDisconnected"
export const ChatAppEventSource = "gabo1208.knative-go-chat-app/source"

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
