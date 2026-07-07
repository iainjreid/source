import { AlertTriangle, Asterisk } from "lucide-react";
import { RefSelector } from "../../components/ref-selector";
import { RepoSelector } from "../../components/repo-selector";
import { TreeNode } from "../../components/tree-node";
import { useLocation } from "../../hooks/use-location";
import { LoadingStatus, useCurrentTree, useCurrentTreeLoading, useSelectedFile, useSelectedFileLoading, useSelectedFilepath } from "../../hooks/use-shared-state";
import { Viewer } from "../viewer/viewer";
import { Link } from "../../components/link";

// ─── Types ───────────────────────────────────────────────────────────────────

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


function toggleNode(nodes: FileNode[], targetPath: string): FileNode[] {
  return nodes.map((node) => {
    if (node.kind === "file") return node;
    if (node.path === targetPath) return { ...node, expanded: !node.expanded };
    return { ...node, children: toggleNode(node.children, targetPath) };
  });
}

// Store the file tree state in the history state, to ensure that accidental page
// changes, or full-page refreshes maintain the users file tree state correctly.
//
// To do this, we only need to store the paths that are expanded, not the entire
// tree.

export function Repo() {
  const currentTreeLoading = useCurrentTreeLoading();
  const currentTree = useCurrentTree();

  const selectedFilepath = useSelectedFilepath();

  function handleToggle(path: string) {
    currentTree.setState((prev) => {
      return toggleNode(prev, path);
    });
  }

  return (
    <div className="flex flex-col h-dvh w-full overflow-hidden bg-background text-foreground">
      {/* ── Navbar ── */}
      <header className="shrink-0 flex items-center gap-0 border-b border-border bg-card">
        <Link className="flex items-center justify-center w-12 h-10 border-r border-border shrink-0" to="/">
          <Asterisk size={16} className="text-accent" strokeWidth={2} />
        </Link>
        <RepoSelector />
        <RefSelector />
      </header>

      {/* ── Body ── */}
      <div className="flex flex-1 min-h-0">
        {/* ── File Tree ── */}
        <aside className="w-56 shrink-0 border-r border-border bg-card flex flex-col min-h-0 overflow-hidden">
          <div className="flex items-center justify-between px-3 py-2 border-b border-border">
            <span className="text-[10px] text-muted-foreground uppercase tracking-widest">Files</span>
          </div>

          <div className="flex-1 overflow-y-auto py-1 relative" style={{ scrollbarWidth: "none" }}>
            {currentTreeLoading.state.status === LoadingStatus.Loading && (
              <div className="absolute inset-0 flex flex-col items-center justify-center gap-2">
                <div className="w-4 h-4 border border-accent/30 border-t-accent rounded-full animate-spin" />
                <span className="text-[10px] text-muted-foreground">Loading…</span>
              </div>
            )}

            {currentTreeLoading.state.status === LoadingStatus.Error && (
              <div className="absolute inset-0 flex flex-col items-center justify-center gap-2 px-4">
                <AlertTriangle size={16} className="text-destructive opacity-70" strokeWidth={1.5} />
                <span className="text-[10px] text-muted-foreground text-center leading-relaxed">
                  {currentTreeLoading.state.message}
                </span>
              </div>
            )}

            {currentTreeLoading.state.status === LoadingStatus.Loaded && (
              currentTree.state?.map((node) => (
                <TreeNode
                  key={node.path}
                  node={node}
                  depth={0}
                  selectedPath={selectedFilepath.state}
                  onSelect={selectedFilepath.setState}
                  onToggle={handleToggle}
                />
              ))
            )}
          </div>
        </aside>

        <Viewer />
      </div>
    </div>
  );
}
