interface FocusToggleButtonProps {
  isFocus: boolean;
  saving: boolean;
  busy: boolean;
  onToggle: () => void;
}

// aria-disabled, not the disabled attribute: a native disabled button is
// forced to blur, which would kick keyboard focus off this button the
// instant its own click starts the request (UX-2's finding).
export function FocusToggleButton({ isFocus, saving, busy, onToggle }: FocusToggleButtonProps) {
  return (
    <button
      type="button"
      onClick={() => {
        if (busy) return;
        onToggle();
      }}
      aria-disabled={busy}
      aria-pressed={isFocus}
      className={`inline-flex shrink-0 items-center gap-1 rounded-full px-3 py-1 text-xs font-medium transition-colors duration-150 ease-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2 ${
        busy ? "opacity-50" : ""
      } ${
        isFocus
          ? "bg-accent text-white"
          : "border border-subtle text-muted hover:bg-accent-tint hover:text-accent-strong"
      }`}
    >
      <span aria-hidden>{isFocus ? "★" : "☆"}</span>
      {saving ? "Saving…" : isFocus ? "Focus" : "Set focus"}
    </button>
  );
}
