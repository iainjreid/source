import { useLocation } from "../hooks/use-location";

export function Link({ className, to, preflight, children }: React.PropsWithChildren<{ className?: string; preflight?: () => void, to: string }>) {
  const location = useLocation();

  function onClick(evt: React.MouseEvent<HTMLAnchorElement>) {
    // Let the browser handle modified clicks and non-left clicks.
    if (
      evt.defaultPrevented ||
      evt.button !== 0 ||
      evt.metaKey ||
      evt.ctrlKey ||
      evt.shiftKey ||
      evt.altKey
    ) {
      return;
    }

    evt.preventDefault();
    preflight?.();
    location.setState(to);
  }

  return (
    <a className={className} href={to} onClick={onClick}>
      {children}
    </a>
  );
}
