import { ChevronsUpDown, Book } from "lucide-react";
import { useState, useRef, useEffect, FC } from "react";
import { Repo } from "../utils/api";

export interface RepoSelectorProps {
  repos: Repo[];
  repo: Repo;
  setRepo: (_: Repo) => void;
}

export const RepoSelector: FC<RepoSelectorProps> = ({
  repos,
  repo,
  setRepo,
}) => {

  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
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
  }, [open]);

  return (
    <div ref={ref} className="relative flex items-center gap-2 h-10 border-r border-border">
      <button
        onClick={() => setOpen((p) => !p)}
        className={`flex items-center gap-1.5 px-4 h-full transition-colors hover:bg-secondary group ${open ? "bg-secondary" : "text-muted-foreground"}`}
      >
        <span className="text-xs text-muted-foreground">{window.location.host}</span>
        <span className="text-xs text-border">/</span>
        <span className="text-xs text-foreground font-medium">{repo.name}</span>
        <ChevronsUpDown size={9} className="ml-0.5 opacity-40" strokeWidth={2} />
      </button>

      {open && (
        <div className="absolute top-full left-[-1px] z-50 w-52 bg-card border border-border shadow-xl shadow-black/60">
          <div className="flex border-b border-border">
            <button className={`flex-1 flex items-center justify-center gap-1.5 py-2 text-[10px] uppercase tracking-widest transition-colors ${true
              ? "text-accent border-b border-accent -mb-px"
              : "text-muted-foreground hover:text-foreground"
              }`}
            >
              <Book size={10} strokeWidth={2} />
              Repo
            </button>
          </div>

          <div className="py-1 max-h-56 overflow-y-auto" style={{ scrollbarWidth: "none" }}>
            {repos.map((r) => (
              <button
                key={r.name}
                onClick={() => { setRepo(r); setOpen(false); }}
                className={`w-full flex items-center justify-between px-3 py-1.5 text-xs text-left transition-colors hover:bg-secondary ${r === repo ? "text-accent" : "text-muted-foreground hover:text-foreground"
                  }`}
              >
                <span>{r.name}</span>
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
