import { useEffect, useRef, useCallback } from "react";

interface UseIntersectionObserverOptions {
  root?: Element | null;
  rootMargin?: string;
  threshold?: number | number[];
}

interface UseIntersectionObserverReturn<T extends HTMLElement = HTMLElement> {
  ref: React.RefObject<T | null>;
}

export const useIntersectionObserver = <T extends HTMLElement = HTMLElement>(
  callback: (entries: IntersectionObserverEntry[]) => void,
  options: UseIntersectionObserverOptions = {}
): UseIntersectionObserverReturn<T> => {
  const ref = useRef<T | null>(null);

  const handleObserver = useCallback(
    (entries: IntersectionObserverEntry[]) => {
      callback(entries);
    },
    [callback]
  );

  useEffect(() => {
    const node = ref.current;
    if (!node) return;

    const observer = new window.IntersectionObserver(handleObserver, {
      root: options.root ?? null,
      rootMargin: options.rootMargin ?? "0px",
      threshold: options.threshold ?? 1.0,
    });

    observer.observe(node);

    return () => {
      observer.unobserve(node);
    };
  }, [handleObserver, options.root, options.rootMargin, options.threshold]);

  return { ref };
};
