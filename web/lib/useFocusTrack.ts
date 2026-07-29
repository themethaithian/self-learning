"use client";

import { useCallback, useState } from "react";
import { setFocusTrack as putFocusTrack } from "./api";

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
  const [focusTrack, setFocusTrackValue] = useState<string | null>(null);
  const [pendingTrack, setPendingTrack] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const setFocus = useCallback(
    async (track: string | null) => {
      const previous = focusTrack;
      setFocusTrackValue(track);
      setPendingTrack(track ?? previous);
      setError(null);
      try {
        await putFocusTrack(track);
      } catch {
        setFocusTrackValue(previous);
        setError("Could not save focus track — please try again.");
      } finally {
        setPendingTrack(null);
      }
    },
    [focusTrack],
  );

  const dismiss = useCallback(() => setError(null), []);

  return { focusTrack, setFocusTrackValue, pendingTrack, error, setFocus, dismiss };
}
