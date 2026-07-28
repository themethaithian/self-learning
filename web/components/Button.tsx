import Link from "next/link";
import type { ButtonHTMLAttributes } from "react";
import type { ComponentProps } from "react";

type Variant = "primary" | "ghost" | "danger";

const VARIANT_CLASSES: Record<Variant, string> = {
  primary: "bg-accent text-white hover:bg-accent-hover",
  ghost: "bg-transparent text-body border border-subtle hover:bg-accent-tint",
  danger: "bg-transparent text-danger border border-subtle hover:bg-danger/10",
};

const BASE_CLASSES =
  "inline-flex items-center justify-center gap-2 rounded-xl px-4 py-2 text-sm font-medium transition-colors duration-150 ease-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
}

export function Button({ variant = "primary", className = "", ...props }: ButtonProps) {
  return (
    <button
      className={`${BASE_CLASSES} disabled:pointer-events-none disabled:opacity-50 ${VARIANT_CLASSES[variant]} ${className}`}
      {...props}
    />
  );
}

interface LinkButtonProps extends ComponentProps<typeof Link> {
  variant?: Variant;
}

// Same button look as <Button>, for actions that are really navigation
// (internal links need a real <a> — middle-click/new-tab, no extra history
// entries — not an onClick + router.push).
export function LinkButton({ variant = "ghost", className = "", ...props }: LinkButtonProps) {
  return <Link className={`${BASE_CLASSES} ${VARIANT_CLASSES[variant]} ${className}`} {...props} />;
}
