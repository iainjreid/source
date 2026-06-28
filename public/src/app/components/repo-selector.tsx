import { GitBranch, Tag, ChevronsUpDown, AlertTriangle, Check } from "lucide-react";
import { useState, useRef, useEffect, FC } from "react";

export interface RepoSelectorProps{
  refKind: RefKind;
  refValue: { name: string, hash: string };
  refsState: RefsState;
  onChange: (kind: RefKind, value: { name: string, hash: string }) => void;
}

export const RepoSelector: FC<RepoSelectorProps> = ({
  refValue,
  refsState,
  onChange,
}) => {
  const [open, setOpen] = useState(false);
  const [tab, setTab] = useState<RefKind>(refKind);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    setTab(refKind);
    function onMouse(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    }
    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") setOpen(false);
    }
    document.addEventListener("mousedown", onMouse);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onMouse);
      document.removeEventListener("keydown", onKey);
    };
  }, [open, refKind]);

  const Icon = refKind === "branch" ? GitBranch : Tag;
  const list = refsState.status === "loaded"
    ? (tab === "branch" ? refsState.branches : refsState.tags)
    : [];

  return (
    <div ref={ref} className="relative flex items-center h-10 border-r border-border">
      <button
        onClick={() => setOpen((p) => !p)}
        className={`flex items-center gap-1.5 px-4 h-full transition-colors hover:bg-secondary group ${open ? "bg-secondary" : "text-muted-foreground"}`}
      >
        <Icon size={11} strokeWidth={2} className={open ? "text-accent" : ""} />
        <span className={`text-xs transition-colors ${open ? "text-foreground" : "group-hover:text-foreground"}`}>
          {refValue?.name}
        </span>
        <ChevronsUpDown size={9} className="ml-0.5 opacity-40" strokeWidth={2} />
      </button>

      {open && (
        <div className="absolute top-full left-0 mt-px z-50 w-52 bg-card border border-border shadow-xl shadow-black/60">
          <div className="flex border-b border-border">
              <button
                key={t}
                onClick={() => setTab(t)}
                className={`flex-1 flex items-center justify-center gap-1.5 py-2 text-[10px] uppercase tracking-widest transition-colors ${tab === t
                  ? "text-accent border-b border-accent -mb-px"
                  : "text-muted-foreground hover:text-foreground"
                  }`}
              >
                {t === "branch" ? <GitBranch size={10} strokeWidth={2} /> : <Tag size={10} strokeWidth={2} />}
                {t}
              </button>
          </div>
          <div className="py-1 max-h-56 overflow-y-auto" style={{ scrollbarWidth: "none" }}>
            {refsState.status === "loading" && (
              <div className="flex items-center justify-center py-4">
                <div className="w-3.5 h-3.5 border border-accent/30 border-t-accent rounded-full animate-spin" />
              </div>
            )}
            {refsState.status === "error" && (
              <div className="flex items-center justify-center gap-1.5 py-4 px-3">
                <AlertTriangle size={11} className="text-destructive shrink-0" strokeWidth={1.5} />
                <span className="text-[10px] text-muted-foreground truncate">{refsState.message}</span>
              </div>
            )}
            {refsState.status === "loaded" && list.map((opt: { name: string, hash: string }) => {
              const isActive = refKind === tab && refValue === opt;
              return (
                <button
                  key={opt.name}
                  onClick={() => { onChange(tab, opt); setOpen(false); }}
                  className={`w-full flex items-center justify-between px-3 py-1.5 text-xs text-left transition-colors hover:bg-secondary ${isActive ? "text-accent" : "text-muted-foreground hover:text-foreground"
                    }`}
                >
                  <span>{opt.name}</span>
                  {isActive && <Check size={10} strokeWidth={2.5} />}
                </button>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}
