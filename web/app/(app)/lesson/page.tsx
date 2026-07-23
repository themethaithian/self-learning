"use client";

import { Suspense, useCallback, useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { getLesson, NotFoundError, UnauthorizedError, type Lesson } from "@/lib/api";
import { LessonBody } from "@/components/LessonBody";
import { RecallCheckCard } from "@/components/RecallCheckCard";
import { EmptyState } from "@/components/EmptyState";
import { Button } from "@/components/Button";
import { LessonSkeleton } from "@/components/Skeleton";
import { BookIcon, WarningIcon } from "@/components/icons";

type ViewState =
  | { status: "loading" }
  | { status: "error"; message: string }
  | { status: "not-found" }
  | { status: "success"; lesson: Lesson };

function BackLink() {
  return (
    <Link
      href="/read"
      className="inline-flex items-center gap-1 rounded-lg text-sm font-medium text-muted transition-colors duration-150 ease-out hover:text-accent focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2"
    >
      ← Back to curriculum tree
    </Link>
  );
}

function LessonView() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const topicSlug = searchParams.get("topic") ?? "";
  const conceptSlug = searchParams.get("concept") ?? "";
  const [state, setState] = useState<ViewState>({ status: "loading" });

  const load = useCallback(async () => {
    if (!topicSlug || !conceptSlug) {
      setState({ status: "error", message: "Missing topic or concept in the URL." });
      return;
    }
    setState({ status: "loading" });
    try {
      const lesson = await getLesson(topicSlug, conceptSlug);
      setState({ status: "success", lesson });
    } catch (err) {
      if (err instanceof UnauthorizedError) {
        router.replace("/token");
        return;
      }
      if (err instanceof NotFoundError) {
        setState({ status: "not-found" });
        return;
      }
      setState({
        status: "error",
        message: err instanceof Error ? err.message : "Something went wrong",
      });
    }
  }, [topicSlug, conceptSlug, router]);

  useEffect(() => {
    load();
  }, [load]);

  if (state.status === "loading") {
    return (
      <div className="space-y-6">
        <BackLink />
        <LessonSkeleton />
      </div>
    );
  }

  if (state.status === "error") {
    return (
      <div className="space-y-6">
        <BackLink />
        <EmptyState
          icon={<WarningIcon />}
          message={`Could not load this lesson: ${state.message}`}
          action={<Button onClick={load}>Retry</Button>}
        />
      </div>
    );
  }

  if (state.status === "not-found") {
    return (
      <div className="space-y-6">
        <BackLink />
        <EmptyState
          icon={<BookIcon />}
          message="No lesson yet for this concept. Most of the curriculum doesn't have a lesson written yet — check back later."
          action={
            <Button variant="ghost" onClick={() => router.push("/read")}>
              Back to curriculum tree
            </Button>
          }
        />
      </div>
    );
  }

  const { lesson } = state;
  return (
    <div className="space-y-8">
      <BackLink />

      <div className="space-y-2">
        <h1 className="text-xl font-semibold text-heading">{lesson.title_en}</h1>
        <span className="inline-flex items-center rounded-full bg-accent-tint px-3 py-1 text-xs font-medium text-accent-strong">
          {lesson.est_minutes} min read
        </span>
      </div>

      <LessonBody markdown={lesson.body_md} />

      <section className="max-w-[68ch] space-y-3">
        <h2 className="text-xs font-medium uppercase tracking-wide text-faint">References</h2>
        <ul className="space-y-3">
          {lesson.references.map((ref) => (
            <li key={ref.title} className="rounded-xl border border-subtle bg-surface p-4">
              <p className="text-sm font-medium text-heading">{ref.title}</p>
              <p className="text-sm text-muted">{ref.source}</p>
              <p className="mt-1 font-thai text-sm leading-[1.8] text-body">{ref.why}</p>
            </li>
          ))}
        </ul>
      </section>

      <section className="max-w-[68ch] space-y-3">
        <h2 className="text-xs font-medium uppercase tracking-wide text-faint">Recall check</h2>
        <div className="space-y-3">
          {lesson.recall_checks.map((check, index) => (
            <RecallCheckCard key={check.position} check={check} index={index} />
          ))}
        </div>
      </section>
    </div>
  );
}

export default function LessonPage() {
  return (
    <Suspense
      fallback={
        <div className="space-y-6">
          <BackLink />
          <LessonSkeleton />
        </div>
      }
    >
      <LessonView />
    </Suspense>
  );
}
