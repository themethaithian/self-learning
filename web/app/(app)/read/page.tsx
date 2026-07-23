"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getCurriculum, UnauthorizedError, type Track } from "@/lib/api";
import { CurriculumTree } from "@/components/CurriculumTree";
import { CurriculumTreeSkeleton } from "@/components/Skeleton";
import { EmptyState } from "@/components/EmptyState";
import { Button } from "@/components/Button";
import { BookIcon, WarningIcon } from "@/components/icons";

type ViewState =
  | { status: "loading" }
  | { status: "error"; message: string }
  | { status: "empty" }
  | { status: "success"; tracks: Track[] };

export default function ReadPage() {
  const router = useRouter();
  const [state, setState] = useState<ViewState>({ status: "loading" });

  const load = useCallback(async () => {
    setState({ status: "loading" });
    try {
      const data = await getCurriculum();
      setState(data.tracks.length === 0 ? { status: "empty" } : { status: "success", tracks: data.tracks });
    } catch (err) {
      if (err instanceof UnauthorizedError) {
        router.replace("/token");
        return;
      }
      setState({
        status: "error",
        message: err instanceof Error ? err.message : "Something went wrong",
      });
    }
  }, [router]);

  useEffect(() => {
    load();
  }, [load]);

  if (state.status === "loading") {
    return <CurriculumTreeSkeleton />;
  }

  if (state.status === "error") {
    return (
      <EmptyState
        icon={<WarningIcon />}
        message={`Could not load the curriculum: ${state.message}`}
        action={<Button onClick={load}>Retry</Button>}
      />
    );
  }

  if (state.status === "empty") {
    return <EmptyState icon={<BookIcon />} message="No curriculum has been imported yet." />;
  }

  return <CurriculumTree tracks={state.tracks} />;
}
