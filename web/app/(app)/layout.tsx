"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import type { ReactNode } from "react";
import { AppShell } from "@/components/AppShell";
import { getToken } from "@/lib/token";

export default function AppGroupLayout({ children }: { children: ReactNode }) {
  const router = useRouter();
  const [ready, setReady] = useState(false);

  useEffect(() => {
    if (getToken()) {
      setReady(true);
    } else {
      router.replace("/token");
    }
  }, [router]);

  if (!ready) {
    return null;
  }

  return <AppShell>{children}</AppShell>;
}
