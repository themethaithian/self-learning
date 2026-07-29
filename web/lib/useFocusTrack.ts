import { useCallback, useState } from "react";
import { useRouter } from "next/navigation";
import { setFocusTrack as putFocusTrack, UnauthorizedError } from "./api";

export interface UseFocusTrackResult {
  focusTrack: string | null;
  setFocusTrackValue: (track: string | null) => void;
  pendingTrack: string | null;
  error: string | null;
  setFocus: (track: string | null) => Promise<void>;
  dismiss: () => void;
}

// Shared by /learn and /today (both render a focus toggle per track) — the
// optimistic-update-with-revert-on-failure logic is exactly what drifts
// silently if copied twice, same reasoning as sortTracksByDisplayOrder.
// setFocusTrackValue is exposed alongside setFocus because the caller loads
// the initial value from GET /api/v1/prefs/focus-track itself; this hook
// only owns the toggle, not the fetch.
export function useFocusTrack(): UseFocusTrackResult {
  const router = useRouter();
  const [focusTrack, setFocusTrackValue] = useState<string | null>(null);
  const [pendingTrack, setPendingTrack] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const setFocus = useCallback(
    async (track: string | null) => {
      // The hook owns this guard rather than trusting a caller-computed
      // `busy` prop — a second toggle while one is in flight would capture
      // a stale `previous` and revert to the wrong value on failure.
      if (pendingTrack !== null) return;

      const previous = focusTrack;
      setFocusTrackValue(track);
      setPendingTrack(track ?? previous);
      setError(null);
      try {
        await putFocusTrack(track);
      } catch (err) {
        // An expired token must still redirect (api.ts's own policy) — a
        // bare catch here would swallow it into "please try again", and
        // retrying would send no Authorization header at all, looping
        // forever without ever reaching /token.
        if (err instanceof UnauthorizedError) {
          router.replace("/token");
          return;
        }
        setFocusTrackValue(previous);
        setError("Could not save focus track — please try again.");
      } finally {
        setPendingTrack(null);
      }
    },
    [focusTrack, pendingTrack, router],
  );

  const dismiss = useCallback(() => setError(null), []);

  return { focusTrack, setFocusTrackValue, pendingTrack, error, setFocus, dismiss };
}
