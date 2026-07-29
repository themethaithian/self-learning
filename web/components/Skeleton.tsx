export function Skeleton({ className = "" }: { className?: string }) {
  return <div className={`animate-pulse rounded-lg bg-subtle/60 ${className}`} />;
}

export function TrackCardsSkeleton() {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2" aria-hidden>
      {Array.from({ length: 4 }).map((_, index) => (
        <div key={index} className="space-y-4 rounded-2xl border border-subtle bg-surface p-6 shadow-sm">
          <div className="flex items-start justify-between gap-2">
            <Skeleton className="h-5 w-40" />
            <Skeleton className="h-6 w-20 rounded-full" />
          </div>
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-9 w-28 rounded-xl" />
        </div>
      ))}
    </div>
  );
}

export function TodaySkeleton() {
  return (
    <div className="space-y-8" aria-hidden>
      <div className="space-y-4 rounded-2xl border border-subtle bg-surface p-6 shadow-sm">
        <Skeleton className="h-3 w-20" />
        <Skeleton className="h-4 w-1/2" />
        <Skeleton className="h-6 w-2/3" />
        <Skeleton className="h-9 w-28 rounded-xl" />
      </div>
      <div className="space-y-2">
        {Array.from({ length: 3 }).map((_, index) => (
          <Skeleton key={index} className="h-16 w-full rounded-xl" />
        ))}
      </div>
    </div>
  );
}

export function LessonSkeleton() {
  return (
    <div className="max-w-[68ch] space-y-6" aria-hidden>
      <Skeleton className="h-4 w-40" />
      <Skeleton className="h-8 w-2/3" />
      <div className="space-y-3">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-5/6" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-2/3" />
      </div>
    </div>
  );
}
