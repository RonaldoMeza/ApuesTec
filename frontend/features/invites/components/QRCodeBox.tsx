"use client";

import { useState, useCallback } from "react";

interface QRCodeBoxProps {
  payload: string;
  label?: string;
}

export function QRCodeBox({ payload, label }: QRCodeBoxProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(payload);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // clipboard not available
    }
  }, [payload]);

  return (
    <div className="rounded-xl border border-border bg-surface p-5">
      {label && (
        <p className="mb-3 text-xs font-medium uppercase tracking-wider text-muted-foreground">
          {label}
        </p>
      )}
      <div className="flex items-center justify-center rounded-lg border border-border bg-surface-muted px-4 py-6">
        <div className="flex flex-col items-center gap-3">
          <div className="flex items-center gap-2">
            <div className="h-3 w-3 rounded-full bg-primary" />
            <div className="h-3 w-3 rounded-full bg-amber-400" />
            <div className="h-3 w-3 rounded-full bg-primary" />
          </div>
          <p className="max-w-full break-all rounded-md bg-black/30 px-3 py-2 font-mono text-xs text-muted-foreground">
            {payload}
          </p>
        </div>
      </div>
      <button
        onClick={handleCopy}
        className="mt-3 w-full rounded-lg border border-border bg-surface-muted px-4 py-2 text-sm text-foreground transition-all hover:bg-primary/10 hover:text-primary"
      >
        {copied ? "Copiado!" : "Copiar enlace"}
      </button>
    </div>
  );
}
