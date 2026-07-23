export function Skeleton({ className = "" }: { className?: string }) {
  return <div className={`animate-pulse rounded-lg bg-subtle/60 ${className}`} />;
}

export function CurriculumTreeSkeleton() {
  return (
    <div className="space-y-8" aria-hidden>
      {Array.from({ length: 3 }).map((_, trackIndex) => (
        <div key={trackIndex} className="space-y-3">
          <Skeleton className="h-5 w-40" />
          <div className="space-y-2 pl-4">
            {Array.from({ length: 3 }).map((_, topicIndex) => (
              <Skeleton key={topicIndex} className="h-11 w-full" />
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
