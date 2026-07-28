import type { SVGProps } from "react";

function BaseIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={1.5}
      strokeLinecap="round"
      strokeLinejoin="round"
      className="h-10 w-10"
      aria-hidden
      {...props}
    />
  );
}

export function BookIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <BaseIcon {...props}>
      <path d="M4 5.5A2.5 2.5 0 0 1 6.5 3H20v15.5a1 1 0 0 1-1 1H6.5A2.5 2.5 0 0 0 4 22z" />
      <path d="M4 19a2.5 2.5 0 0 1 2.5-2.5H19" />
    </BaseIcon>
  );
}

export function WarningIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <BaseIcon {...props}>
      <path d="M10.29 3.86 1.82 18a1 1 0 0 0 .86 1.5h18.64a1 1 0 0 0 .86-1.5L13.71 3.86a1 1 0 0 0-1.72 0Z" />
      <path d="M12 9v4" />
      <path d="M12 17h.01" />
    </BaseIcon>
  );
}

export function CheckIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      viewBox="0 0 20 20"
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
      strokeLinecap="round"
      strokeLinejoin="round"
      className="h-4 w-4 shrink-0"
      aria-hidden
      {...props}
    >
      <path d="M4 10.5l4 4 8-9" />
    </svg>
  );
}

export function ChevronIcon({ open, ...props }: SVGProps<SVGSVGElement> & { open: boolean }) {
  return (
    <svg
      viewBox="0 0 20 20"
      fill="none"
      stroke="currentColor"
      strokeWidth={1.5}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={`h-4 w-4 shrink-0 transition-transform duration-150 ease-out ${open ? "rotate-90" : ""}`}
      aria-hidden
      {...props}
    >
      <path d="M7 4l6 6-6 6" />
    </svg>
  );
}
