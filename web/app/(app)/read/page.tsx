"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

// /read was replaced by /learn (track cards + focus picker). Static export
// has no server-side redirects, so this client component bounces old links
// and bookmarks to the new route on load.
export default function ReadPage() {
  const router = useRouter();

  useEffect(() => {
    router.replace("/learn");
  }, [router]);

  return null;
}
