import { useEffect, useState } from "react";
import { SharedState } from "./use-shared-state";

export function useLocation(): SharedState<string> & { replaceState: (_: string) => void } {
  const [state, setPath] = useState(window.location.pathname);

  useEffect(() => {
    function onPopState() {
      console.log("updating popstate");
      setPath(window.location.pathname);
    };

    window.addEventListener("popstate", onPopState);

    return () => {
      window.removeEventListener("popstate", onPopState);
    };
  }, []);

  const setState = (next: string | ((prev: string) => string)) => {
    const to =
      typeof next === "function"
        ? (next as (prev: string) => string)(state)
        : next;

    if (window.location.pathname === to) {
      return
    }
    console.debug("updating location state", to)
    window.history.pushState(null, "", to);

    var popStateEvent = new PopStateEvent('popstate', { state: null });
    dispatchEvent(popStateEvent);
    console.debug("setting location state", to)
    setPath(to);
  };

  const replaceState = (to: string) => {
    window.history.replaceState(null, "", to);
    var popStateEvent = new PopStateEvent('popstate', { state: null });
    dispatchEvent(popStateEvent);
    console.debug("replacing location state", to)
    setPath(to);
  };

  return { state, setState, replaceState };
}

/**
 * Removes the '+' suffix from a path parameter name.
 *
 * For example:
 *
 *   "ref+"      -> "ref"
 *   "filename"  -> "filename"
 */
type StripSuffix<T extends string> =
  T extends `${infer Name}+` ? Name : T;

/**
 * Produces an object containing a single path parameter.
 *
 * For example:
 *
 *   Param<"repo"> -> { repo: string }
 *   Param<"ref+"> -> { ref: string }
 */
type Param<T extends string> = {
  [K in StripSuffix<T>]: string;
};

/**
 * Produces the object type for the parameters captured by a route.
 *
 * Parameters are introduced by ':' and may be suffixed with '+' to indicate
 * that they may span multiple path segments.
 *
 * For example:
 *
 *   PathParams<"/:repo">
 *   // { repo: string }
 *
 *   PathParams<"/:repo/refs/:ref+/blob/:filename+">
 *   // {
 *   //   repo: string;
 *   //   ref: string;
 *   //   filename: string;
 *   // }
 */
type PathParams<T extends string> =
  T extends `${string}:${infer P}/${infer Rest}`
    ? Param<P> & PathParams<`/${Rest}`>
    : T extends `${string}:${infer P}`
      ? Param<P>
      : {};

export function match<Pattern extends string>(pattern: Pattern, path: string): Partial<PathParams<Pattern>> {
  const parts = pattern.split("/").filter(Boolean);

  let regex = "^";

  for (let i = 0; i < parts.length; i++) {
    const part = parts[i];

    if (part.startsWith(":")) {
      const greedy = part.endsWith("+");
      const name = part.slice(1, greedy ? -1 : undefined);

      let capture = "[^/]+";

      if (greedy) {
        // If followed by a literal, stop at the first occurrence of it.
        const next = parts[i + 1];

        capture = next && !next.startsWith(":")
          ? ".+?"
          : ".+";
      }

      regex += `(?:/(?<${name}>${capture})`;
    } else {
      regex += `(?:/${part}`;
    }
  }

  regex += ")?".repeat(parts.length);
  regex += "$";

  return (path.match(regex)?.groups ?? {}) as Partial<PathParams<Pattern>>;
}
