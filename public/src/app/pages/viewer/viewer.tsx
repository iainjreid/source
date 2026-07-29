import { Copy, File } from "lucide-react";
import { Code } from "../../components/code";
import { LoadingStatus, useSelectedFile, useSelectedFileLoading } from "../../hooks/use-shared-state";
import { formatBytes } from "../../utils/format-bytes";

export function Viewer() {
  const selectedFile = useSelectedFile();
  const selectedFileLoading = useSelectedFileLoading();

  return <>
    {/* ── File Viewer ── */}
    <main className="flex-1 flex flex-col min-h-0 min-w-0">
      {/* File header */}
      <div className="shrink-0 flex items-center justify-between px-4 py-2 border-b border-border bg-card min-h-[37px]">
        {selectedFileLoading.state.status === LoadingStatus.Idle && (
          <span className="text-[10px] text-muted-foreground">Select a file to view</span>
        )}

        {selectedFileLoading.state.status === LoadingStatus.Loading && (
          <div className="flex items-center gap-3">
            <div className="w-24 h-2 bg-secondary rounded animate-pulse" />
            <div className="w-16 h-2 bg-secondary rounded animate-pulse" />
            <div className="w-36 h-2 bg-secondary rounded animate-pulse" />
          </div>
        )}

        {selectedFileLoading.state.status === LoadingStatus.Loaded && (
          <>
            <div className="flex items-center gap-3 min-w-0">
              <div className="flex items-center gap-1.5">
                <File size={12} className="text-accent shrink-0" strokeWidth={1.5} />
                <span className="text-xs font-medium text-foreground truncate">{selectedFile.state!.fileName}</span>
              </div>
              <span className="text-[10px] text-border">·</span>
              <span className="text-[10px] text-muted-foreground">
                {selectedFile.state!.lineCount.toLocaleString()} lines
              </span>
              <span className="text-[10px] text-border">·</span>
              <span className="text-[10px] text-muted-foreground">{formatBytes(selectedFile.state!.fileBytes)}</span>
              <span className="hidden sm:inline text-[10px] text-border">·</span>
              <span
                className="hidden sm:inline text-[10px] text-muted-foreground truncate max-w-xs"
                title={selectedFile.state!.lastCommitHash}
              >
                {selectedFile.state!.lastCommitMsg}
              </span>
            </div>
            <div className="flex items-center gap-2 shrink-0 ml-4">
              <span className="text-[10px] text-muted-foreground hidden md:inline">
                by <span className="text-accent">{selectedFile.state!.lastCommitAuthor}</span>
              </span>
              <button className="flex items-center gap-1.5 px-2.5 py-1 border border-border hover:border-accent/40 hover:bg-secondary rounded transition-all text-muted-foreground hover:text-foreground">
                <Copy size={11} strokeWidth={1.5} />
                <span className="text-[10px]">Raw</span>
              </button>
            </div>
          </>
        )}
      </div>

      {/* Code area */}
      <div
        className="flex-1 overflow-auto overscroll-none relative"
        style={{ scrollbarWidth: "thin", scrollbarColor: "rgba(61,255,160,0.15) transparent" }}
      >
        {selectedFileLoading.state.status === LoadingStatus.Idle && (
          <div className="absolute inset-0 flex flex-col items-center justify-center gap-2 select-none">
            <File size={28} strokeWidth={1} className="text-muted-foreground opacity-20" />
            <span className="text-xs text-muted-foreground opacity-60">Select a file from the tree</span>
          </div>
        )}

        {selectedFileLoading.state.status === LoadingStatus.Loading && (
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="flex flex-col items-center gap-3">
              <div className="w-5 h-5 border border-accent/30 border-t-accent rounded-full animate-spin" />
              <span className="text-[10px] text-muted-foreground">Loading…</span>
            </div>
          </div>
        )}

        {selectedFileLoading.state.status === LoadingStatus.Error && (
          <div className="absolute inset-0 flex flex-col items-center justify-center gap-2 px-8">
            <span className="text-xs text-destructive font-medium">Failed to load file</span>
            <span className="text-[11px] text-muted-foreground text-center leading-relaxed">
              {selectedFileLoading.state.message}
            </span>
            <button
              onClick={() => {
                const path = params.filepath!;
                setSelectedPath("");
                setTimeout(() => setSelectedPath(path), 0);
              }}
              className="mt-2 text-[10px] px-3 py-1 border border-border hover:border-accent/40 hover:bg-secondary text-muted-foreground hover:text-foreground transition-all"
            >
              Retry
            </button>
          </div>
        )}

        {selectedFileLoading.state.status === LoadingStatus.Loaded && (
          <Code lines={selectedFile.state!.fileLines || []} filePath={selectedFile.state!.fileName} />
        )}
      </div>
    </main>
  </>;
}
