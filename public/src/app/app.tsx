import { useState, useEffect, useRef } from "react";
import { GitBranch, Star, GitFork, Eye, File, Check, Circle, Github, Globe, Clock, Tag, ChevronsUpDown, Copy, AlertTriangle, GitMerge } from "lucide-react";
import { Code } from "./components/code";
import { GetFile } from "./utils/api";
import { TreeNode } from "./components/tree-node";
import { Bluesky } from "./components/bluesky";

// ─── Constants ───────────────────────────────────────────────────────────────

const REPO = "source";

// ─── Types ───────────────────────────────────────────────────────────────────

type RefKind = "branch" | "tag";

type ApiRefs = {
  branches: Record<string, string>;
  tags: Record<string, string>;
};

type RefsState =
  | { status: "loading" }
  | { status: "loaded"; branches: { name: string, hash: string }[]; tags: { name: string, hash: string }[] }
  | { status: "error"; message: string };

type ApiTreeEntry = {
  name: string;
  path: string;
  isFile: boolean;
  children: ApiTree;
};

type ApiTree = {
  dirs: Record<string, ApiTreeEntry>;
  files: Record<string, ApiTreeEntry>;
};

type FileNode =
  | { kind: "dir"; name: string; path: string; children: FileNode[]; expanded: boolean }
  | { kind: "file"; name: string; path: string };

type TreeState =
  | { status: "loading" }
  | { status: "loaded"; nodes: FileNode[] }
  | { status: "error"; message: string };

type FileResponse = {
  fileName: string;
  fileBytes: number;
  fileLines: string[];
  lineCount: number;
  lastCommitHash: string;
  lastCommitMsg: string;
  lastCommitAuthor: string;
  timeElapsed: number;
};

type FileState =
  | { status: "idle" }
  | { status: "loading" }
  | { status: "loaded"; data: FileResponse }
  | { status: "error"; message: string };

// ─── Helpers ─────────────────────────────────────────────────────────────────

function formatBytes(n: number): string {
  if (n < 1024) return `${n} B`;
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`;
  return `${(n / (1024 * 1024)).toFixed(1)} MB`;
}

function convertTree(apiTree: ApiTree, selectedPath: string): FileNode[] {
  const dirs: FileNode[] = Object.values(apiTree.dirs).map((entry) => ({
    kind: "dir",
    name: entry.name,
    path: entry.path,
    expanded: selectedPath.startsWith(entry.path),
    children: convertTree(entry.children, selectedPath),
  }));
  const files: FileNode[] = Object.values(apiTree.files).map((entry) => ({
    kind: "file",
    name: entry.name,
    path: entry.path,
  }));
  return [...dirs, ...files];
}

function toggleNode(nodes: FileNode[], targetPath: string): FileNode[] {
  return nodes.map((node) => {
    if (node.kind === "file") return node;
    if (node.path === targetPath) return { ...node, expanded: !node.expanded };
    return { ...node, children: toggleNode(node.children, targetPath) };
  });
}

// ─── RefSelector ─────────────────────────────────────────────────────────────

function RefSelector({
  refKind,
  refValue,
  refsState,
  onChange,
}: {
  refKind: RefKind;
  refValue: { name: string, hash: string };
  refsState: RefsState;
  onChange: (kind: RefKind, value: { name: string, hash: string }) => void;
}) {
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
            {(["branch", "tag"] as RefKind[]).map((t) => (
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
            ))}
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

// ─── App ─────────────────────────────────────────────────────────────────────

export function App() {
  const path = window.location.pathname.slice(1);

  const [refsState, setRefsState] = useState<RefsState>({ status: "loading" });
  const [treeState, setTreeState] = useState<TreeState>({ status: "loading" });
  const [selectedPath, setSelectedPath] = useState<string>(path);
  const [fileState, setFileState] = useState<FileState>({ status: "idle" });
  const [refKind, setRefKind] = useState<RefKind>("branch");
  const [refValue, setRefValue] = useState<{ name: string, hash: string }>();

  useEffect(() => {
    const url = "/" + selectedPath;
    window.history.pushState({ selectedPath }, "", url);
  }, [selectedPath]);

  useEffect(() => {
    const onPopState = () => {
      const path = window.location.pathname.slice(1);
      setSelectedPath(decodeURIComponent(path));
    };

    window.addEventListener("popstate", onPopState);
    return () => window.removeEventListener("popstate", onPopState);
  }, []);

  function handleRefChange(kind: RefKind, value: { name: string, hash: string }) {
    setRefKind(kind);
    setRefValue(value);
    setSelectedPath("");
    setFileState({ status: "idle" });
  }

  function handleToggle(path: string) {
    setTreeState((prev) => {
      if (prev.status !== "loaded") return prev;
      return { ...prev, nodes: toggleNode(prev.nodes, path) };
    });
  }

  // Fetch branches and tags once on mount
  useEffect(() => {
    const controller = new AbortController();
    fetch(`http://localhost:8080/019f0abd-ffd8-7518-80ee-138dd420badc/refs`, { signal: controller.signal })
      .then(async (res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status} — ${res.statusText}`);
        const data: ApiRefs = await res.json();
        const branches = Object.entries(data.branches)
          .map((k) => ({ name: k[0].slice("refs/heads/".length), hash: k[1] }))
        const tags = Object.entries(data.tags)
          .map((k) => ({ name: k[0].slice("refs/tags/".length), hash: k[1] }))
          .sort((a, b) => b.name.localeCompare(a.name));
        setRefsState({ status: "loaded", branches, tags });
        if (branches.length > 0) {
          setRefKind("branch");
          setRefValue(branches[0]);
        } else if (tags.length > 0) {
          setRefKind("tag");
          setRefValue(tags[0]);
        }
      })
      .catch((err: unknown) => {
        if ((err as Error).name === "AbortError") return;
        setRefsState({ status: "error", message: err instanceof Error ? err.message : String(err) });
      });
    return () => controller.abort();
  }, []);

  // Fetch file tree whenever ref changes
  useEffect(() => {
    if (!refValue) return;
    const controller = new AbortController();
    const url = `http://localhost:8080/${REPO}/tree/${encodeURIComponent(refValue.hash)}`;

    setTreeState({ status: "loading" });

    const timeout = setTimeout(() => controller.abort(), 15_000);

    fetch(url, { signal: controller.signal })
      .then(async (res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status} — ${res.statusText}`);
        const data: ApiTree = await res.json();
        setTreeState({ status: "loaded", nodes: convertTree(data, selectedPath) });
      })
      .catch((err: unknown) => {
        if ((err as Error).name === "AbortError") {
          setTreeState({ status: "error", message: "Request timed out." });
        } else {
          setTreeState({ status: "error", message: err instanceof Error ? err.message : String(err) });
        }
      })
      .finally(() => clearTimeout(timeout));

    return () => {
      controller.abort();
      clearTimeout(timeout);
    };
  }, [refValue]);

  // Fetch file content whenever selected path or ref changes
  useEffect(() => {
    if (!selectedPath) return;

    window.history.pushState
    setFileState({ status: "loading" });

    const { destructor, promise } = GetFile("source" /*" 019f0abd-ffd8-7518-80ee-138dd420badc"*/, refValue?.hash!, selectedPath)

    promise
      .then((file) => {
        setFileState({ status: "loaded", data: file });
      })
      .catch((err: unknown) => {
        setFileState({ status: "error", message: err instanceof Error ? err.message : String(err) });
      });

    return destructor;
  }, [selectedPath, refValue]);

  const treeNodes = treeState.status === "loaded" ? treeState.nodes : [];

  return (
    <div
      className="flex flex-col h-dvh overflow-hidden bg-background text-foreground"
      style={{ fontFamily: "'JetBrains Mono', monospace" }}
    >
      {/* ── Navbar ── */}
      <header className="shrink-0 flex items-center gap-0 border-b border-border bg-card">
        <div className="flex items-center justify-center w-12 h-10 border-r border-border shrink-0">
          <GitMerge size={16} className="text-accent" strokeWidth={2} />
        </div>

        <div className="flex items-center gap-2 px-4 h-10 border-r border-border">
          <span className="text-xs text-muted-foreground">{window.location.host}</span>
          <span className="text-xs text-border">/</span>
          <span className="text-xs text-foreground font-medium">{REPO}</span>
        </div>

        <RefSelector refKind={refKind} refValue={refValue!} refsState={refsState} onChange={handleRefChange} />
      </header>

      {/* ── Body ── */}
      <div className="flex flex-1 min-h-0">
        {/* ── File Tree ── */}
        <aside className="w-56 shrink-0 border-r border-border bg-card flex flex-col min-h-0 overflow-hidden">
          <div className="flex items-center justify-between px-3 py-2 border-b border-border">
            <span className="text-[10px] text-muted-foreground uppercase tracking-widest">Files</span>
            <span className="text-[10px] text-muted-foreground">{refValue?.name}</span>
          </div>

          <div className="flex-1 overflow-y-auto py-1 relative" style={{ scrollbarWidth: "none" }}>
            {treeState.status === "loading" && (
              <div className="absolute inset-0 flex flex-col items-center justify-center gap-2">
                <div className="w-4 h-4 border border-accent/30 border-t-accent rounded-full animate-spin" />
                <span className="text-[10px] text-muted-foreground">Loading…</span>
              </div>
            )}

            {treeState.status === "error" && (
              <div className="absolute inset-0 flex flex-col items-center justify-center gap-2 px-4">
                <AlertTriangle size={16} className="text-destructive opacity-70" strokeWidth={1.5} />
                <span className="text-[10px] text-muted-foreground text-center leading-relaxed">
                  {treeState.message}
                </span>
              </div>
            )}

            {treeState.status === "loaded" && treeNodes.map((node) => (
              <TreeNode
                key={node.path}
                node={node}
                depth={0}
                selectedPath={selectedPath}
                onSelect={setSelectedPath}
                onToggle={handleToggle}
              />
            ))}
          </div>
        </aside>

        {/* ── File Viewer ── */}
        <main className="flex-1 flex flex-col min-h-0 min-w-0">
          {/* File header */}
          <div className="shrink-0 flex items-center justify-between px-4 py-2 border-b border-border bg-card min-h-[37px]">
            {fileState.status === "loaded" ? (
              <>
                <div className="flex items-center gap-3 min-w-0">
                  <div className="flex items-center gap-1.5">
                    <File size={12} className="text-accent shrink-0" strokeWidth={1.5} />
                    <span className="text-xs font-medium text-foreground truncate">{fileState.data.fileName}</span>
                  </div>
                  <span className="text-[10px] text-border">·</span>
                  <span className="text-[10px] text-muted-foreground">
                    {fileState.data.lineCount.toLocaleString()} lines
                  </span>
                  <span className="text-[10px] text-border">·</span>
                  <span className="text-[10px] text-muted-foreground">{formatBytes(fileState.data.fileBytes)}</span>
                  <span className="hidden sm:inline text-[10px] text-border">·</span>
                  <span
                    className="hidden sm:inline text-[10px] text-muted-foreground truncate max-w-xs"
                    title={fileState.data.lastCommitHash}
                  >
                    {fileState.data.lastCommitMsg}
                  </span>
                </div>
                <div className="flex items-center gap-2 shrink-0 ml-4">
                  <span className="text-[10px] text-muted-foreground hidden md:inline">
                    by <span className="text-accent">{fileState.data.lastCommitAuthor}</span>
                  </span>
                  <button className="flex items-center gap-1.5 px-2.5 py-1 border border-border hover:border-accent/40 hover:bg-secondary rounded transition-all text-muted-foreground hover:text-foreground">
                    <Copy size={11} strokeWidth={1.5} />
                    <span className="text-[10px]">Raw</span>
                  </button>
                </div>
              </>
            ) : fileState.status === "loading" ? (
              <div className="flex items-center gap-3">
                <div className="w-24 h-2 bg-secondary rounded animate-pulse" />
                <div className="w-16 h-2 bg-secondary rounded animate-pulse" />
                <div className="w-36 h-2 bg-secondary rounded animate-pulse" />
              </div>
            ) : (
              <span className="text-[10px] text-muted-foreground">Select a file to view</span>
            )}
          </div>

          {/* Code area */}
          <div
            className="flex-1 overflow-auto overscroll-none relative"
            style={{ scrollbarWidth: "thin", scrollbarColor: "rgba(61,255,160,0.15) transparent" }}
          >
            {fileState.status === "idle" && (
              <div className="absolute inset-0 flex flex-col items-center justify-center gap-2 select-none">
                <File size={28} strokeWidth={1} className="text-muted-foreground opacity-20" />
                <span className="text-xs text-muted-foreground opacity-60">Select a file from the tree</span>
              </div>
            )}

            {fileState.status === "loading" && (
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="flex flex-col items-center gap-3">
                  <div className="w-5 h-5 border border-accent/30 border-t-accent rounded-full animate-spin" />
                  <span className="text-[10px] text-muted-foreground">Loading…</span>
                </div>
              </div>
            )}

            {fileState.status === "error" && (
              <div className="absolute inset-0 flex flex-col items-center justify-center gap-2 px-8">
                <span className="text-xs text-destructive font-medium">Failed to load file</span>
                <span className="text-[11px] text-muted-foreground text-center leading-relaxed">
                  {fileState.message}
                </span>
                <button
                  onClick={() => {
                    const path = selectedPath;
                    setSelectedPath("");
                    setTimeout(() => setSelectedPath(path), 0);
                  }}
                  className="mt-2 text-[10px] px-3 py-1 border border-border hover:border-accent/40 hover:bg-secondary text-muted-foreground hover:text-foreground transition-all"
                >
                  Retry
                </button>
              </div>
            )}

            {fileState.status === "loaded" && (
              <Code lines={fileState.data.fileLines} filePath={fileState.data.fileName} />
            )}
          </div>
        </main>
      </div>

      {/* ── Status Bar ── */}
      <footer className="shrink-0 flex items-center border-t border-border bg-card px-3 h-7 gap-4 text-[10px]">
        <div className="flex items-center gap-1.5">
          <Circle
            size={6}
            className={
              fileState.status === "error" || treeState.status === "error"
                ? "fill-destructive text-destructive"
                : "fill-accent text-accent"
            }
          />
          <span className="text-muted-foreground">
            {fileState.status === "error" || treeState.status === "error"
              ? "ERR"
              : fileState.status === "loading" || treeState.status === "loading"
                ? "…"
                : "OK"}
          </span>
        </div>

        <div className="w-px h-3 bg-border" />

        <div className="flex items-center gap-1.5 text-muted-foreground">
          <Clock size={10} strokeWidth={1.5} />
          <span>
            rendered in{" "}
            <span className="text-foreground/80">
              {fileState.status === "loaded" ? `${fileState.data.timeElapsed.toFixed(1)}ms` : "—"}
            </span>
          </span>
        </div>

        {/* <div className="w-px h-3 bg-border" />

        <div className="flex items-center gap-1.5 text-muted-foreground">
          {refKind === "branch" ? <GitBranch size={10} strokeWidth={1.5} /> : <Tag size={10} strokeWidth={1.5} />}
          <span>{refValue?.name}</span>
        </div> */}

        <div className="flex-1" />

        <div className="flex items-center gap-3">
          {[
            { icon: Github, label: "GitHub", href: "https://github.com/iainjreid" },
            { icon: Bluesky, label: "Bluesky", href: "https://bsky.app/profile/iainjreid.com" },
            { icon: Globe, label: "Website", href: "https://iainjreid.com" },
          ].map(({ icon: Icon, label, href }) => (
            <a
              key={label}
              href={href}
              target="_blank"
              aria-label={label}
              className="flex items-center gap-1.5 text-muted-foreground hover:text-accent transition-colors"
            >
              <Icon size={11} strokeWidth={1.5} />
              <span className="hidden sm:inline">{label}</span>
            </a>
          ))}
        </div>
      </footer>
    </div>
  );
}
