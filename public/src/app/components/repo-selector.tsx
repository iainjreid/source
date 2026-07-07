import { ChevronsUpDown, Book } from "lucide-react";
import { useState, useRef, useEffect } from "react";
import { Link } from "./link";
import { useKnownRepos, useSelectedRepo } from "../hooks/use-shared-state";

export function RepoSelector() {
  const knownRepos = useKnownRepos();
  const selectedRepo = useSelectedRepo();

  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    console.debug("closing repo selector window");
    setOpen(false);
  }, [selectedRepo.state]);

  useEffect(() => {
    if (!open) return;

    function onMouse(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }

    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") {
        setOpen(false);
      }
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
        <span className="text-xs text-foreground font-medium">{selectedRepo.state?.name}</span>
        <ChevronsUpDown size={9} className="ml-0.5 opacity-40" strokeWidth={2} />
      </button>

      {open && (
        <div className="absolute top-full left-[-1px] z-50 w-52 bg-card border border-border shadow-xl shadow-black/60">
          <div className="flex border-b border-border">
            <button className={`flex-1 flex items-center justify-center gap-1.5 py-2 text-[10px] uppercase tracking-widest transition-colors text-accent border-b border-accent -mb-px`}>
              <Book size={10} strokeWidth={2} />
              Repo
            </button>
          </div>

          <div className="py-1 max-h-56 overflow-y-auto" style={{ scrollbarWidth: "none" }}>
            {knownRepos.state.map((r) => (
              <Link className={`w-full flex items-center justify-between px-3 py-1.5 text-xs text-left transition-colors hover:bg-secondary ${r === selectedRepo.state ? "text-accent" : "text-muted-foreground hover:text-foreground"}`} to={`/${r.name}`}>{r.name}</Link>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
