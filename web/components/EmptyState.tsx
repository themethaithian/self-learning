import type { ReactNode } from "react";

interface EmptyStateProps {
  icon: ReactNode;
  message: string;
  action?: ReactNode;
}

export function EmptyState({ icon, message, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center gap-4 py-16 text-center">
      <div className="text-faint" aria-hidden>
        {icon}
      </div>
      <p className="max-w-sm text-sm text-muted">{message}</p>
      {action}
    </div>
  );
}
