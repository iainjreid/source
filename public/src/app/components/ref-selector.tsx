import { Check, ChevronsUpDown, GitBranch, LucideIcon, Tag } from "lucide-react";
import { Dispatch, SetStateAction, useEffect, useRef, useState } from "react";
import { Link } from "./link";
import { useKnownRefs, useSelectedRef, useSelectedRepo } from "../hooks/use-shared-state";
import { RefKind } from "../utils/api";

export function RefSelector() {
  const knownRefs = useKnownRefs();
  const selectedRepo = useSelectedRepo();
  const selectedRef = useSelectedRef();
  const [open, setOpen] = useState(false);
  const [tab, setTab] = useState<RefKind>(RefKind.Branch);
  const ref = useRef<HTMLDivElement>(null);

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
    <div ref={ref} className="relative flex items-center h-10 border-r border-border">
      <button onClick={() => setOpen((p) => !p)} className={`flex items-center gap-1.5 px-4 h-full transition-colors hover:bg-secondary group ${open ? "bg-secondary" : "text-muted-foreground"}`}>
        {selectedRef.state?.kind === RefKind.Branch
          ? <GitBranch size={11} strokeWidth={2} className={open ? "text-accent" : ""} />
          : <Tag size={11} strokeWidth={2} className={open ? "text-accent" : ""} />
        }
        <span className={`text-xs transition-colors ${open ? "text-foreground" : "group-hover:text-foreground"}`}>
          {selectedRef.state?.shortName}
        </span>
        <ChevronsUpDown size={9} className="ml-0.5 opacity-40" strokeWidth={2} />
      </button>

      {open && (
        <div className="absolute top-full left-[-1px] z-50 w-52 bg-card border border-border shadow-xl shadow-black/60">
          <div className="flex border-b border-border">
            <SelectorTabHeader tab={tab} setTab={setTab} icon={GitBranch} name="Branch" tabKey={RefKind.Branch} />
            <SelectorTabHeader tab={tab} setTab={setTab} icon={Tag} name="Tag" tabKey={RefKind.Tag} />
          </div>
          <div className="py-1 max-h-56 overflow-y-auto" style={{ scrollbarWidth: "none" }}>
            {(Object.values((tab === RefKind.Branch ? knownRefs.state.branches : knownRefs.state.tags) ?? {})).map((ref) => {
              const isActive = selectedRef.state?.hash === ref.hash
              return (
                <Link className={`w-full flex items-center justify-between px-3 py-1.5 text-xs text-left transition-colors hover:bg-secondary ${isActive ? "text-accent" : "text-muted-foreground hover:text-foreground"}`} to={`/${selectedRepo.state?.name}/refs/${ref.shortName}`} preflight={() => setOpen(false)}>
                  <span>{ref.shortName}</span>
                  {isActive && <Check size={10} strokeWidth={2.5} />}
                </Link>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}

interface SelectorTabHeaderProps {
  tab: RefKind;
  setTab: Dispatch<SetStateAction<RefKind>>;
  icon: LucideIcon;
  name: string;
  tabKey: RefKind;
}

function SelectorTabHeader(props: SelectorTabHeaderProps) {
  return (
    <button onClick={() => props.setTab(props.tabKey)}
      className={`flex-1 flex items-center justify-center gap-1.5 py-2 text-[10px] uppercase tracking-widest transition-colors ${props.tab === props.tabKey
        ? "text-accent border-b border-accent -mb-px"
        : "text-muted-foreground hover:text-foreground"
        }`}
    >
      <props.icon size={10} strokeWidth={2} />
      {props.name}
    </button>
  );
}

interface SelectorTabListProps<T> {
  list: T[];
  isActive: (_: T) => boolean;
}

function SelectorTabList(props: SelectorTabListProps<T>) {
  return <>
    {(Object.values((tab === RefKind.Branch ? knownRefs.state.branches : knownRefs.state.tags) ?? {})).map((ref) => {
      const isActive = selectedRef.state?.hash === ref.hash
      return (
        <Link className={`w-full flex items-center justify-between px-3 py-1.5 text-xs text-left transition-colors hover:bg-secondary ${isActive ? "text-accent" : "text-muted-foreground hover:text-foreground"}`} to={`/${selectedRepo.state?.name}/${ref.name}`}>
          <span>{ref.name}</span>
          {isActive && <Check size={10} strokeWidth={2.5} />}
        </Link>
      );
    })}
  </>
}
