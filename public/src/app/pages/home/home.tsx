import { useEffect } from "react";
import { AlertTriangle, ArrowRight, Asterisk, Circle, Clock, Globe } from "lucide-react";
import { Link } from "../../components/link";
import { match, useLocation } from "../../hooks/use-location";
import { useKnownReposLoading, useKnownRepos, useSelectedRepo, useSelectedRef, useKnownRefs, LoadingStatus, useKnownRefsLoading, useCurrentTree, FileNode, useCurrentTreeLoading, useSelectedFilepath, useSelectedFileLoading, useSelectedFile } from "../../hooks/use-shared-state";
import { Repo } from "../repo/repo";
import { GetRepos, Ref, GetRefs, RefKind, GetTree, Tree, GetFile } from "../../utils/api";
import { Github } from "../../components/github";
import { Bluesky } from "../../components/bluesky";

export function Home() {
  const knownRepos = useKnownRepos();
  const knownReposLoading = useKnownReposLoading();
  const selectedRepo = useSelectedRepo();

  const knownRefs = useKnownRefs();
  const knownRefsLoading = useKnownRefsLoading();
  const selectedRef = useSelectedRef();

  const currentTree = useCurrentTree();
  const currentTreeLoading = useCurrentTreeLoading();

  const selectedFilepath = useSelectedFilepath();

  const selectedFile = useSelectedFile();
  const selectedFileLoading = useSelectedFileLoading();

  const location = useLocation();

  useEffect(() => {
    knownReposLoading.setLoading();
    GetRepos()
      .then((repos) => {
        selectedRepo.setState(undefined);
        knownRepos.setState(Object.values(repos));
        knownReposLoading.setLoaded();
      })
      .catch((err) => {
        knownRefsLoading.setError(err);
      });
  }, []);

  useEffect(() => {
    if (!selectedRepo.state) {
      return;
    }

    knownRefsLoading.setLoading();
    GetRefs(selectedRepo.state)
      .then((refs) => {
        const branches: Record<string, Ref> = {};
        const tags: Record<string, Ref> = {};

        for (const branch in refs.branches) {
          const branchName = branch.slice("refs/heads/".length)
          branches[branchName] = {
            name: branch,
            shortName: branchName,
            hash: refs.branches[branch],
            kind: RefKind.Branch,
          }
        }

        for (const tag in refs.tags) {
          const tagName = tag.slice("refs/tags/".length)
          tags[tagName] = {
            name: tag,
            shortName: tagName,
            hash: refs.tags[tag],
            kind: RefKind.Tag,
          }
        }

        selectedRef.setState(undefined);
        knownRefs.setState({
          branches,
          tags,
        });
        knownRefsLoading.setLoaded();
      })
      .catch((err) => {
        knownRefsLoading.setError(err);
      });
  }, [selectedRepo.state]);

  useEffect(() => {
    if (!selectedRepo.state || !selectedRef.state) {
      return;
    }

    currentTreeLoading.setLoading();
    GetTree(selectedRepo.state, selectedRef.state)
      .then((tree) => {
        currentTree.setState(convertTree(tree))
        currentTreeLoading.setLoaded();
      })
      .catch((err) => {
        currentTreeLoading.setError(err);
      })
  }, [selectedRef.state]);

  useEffect(() => {
    if (!selectedRepo.state || !selectedRef.state || !selectedFilepath.state) {
      return;
    }

    selectedFileLoading.setLoading();
    GetFile(selectedRepo.state, selectedRef.state, selectedFilepath.state)
      .then((file) => {
        selectedFile.setState(file)
        selectedFileLoading.setLoaded();
      }).catch((err) => {
        selectedFileLoading.setError(err);
      });
  }, [selectedFilepath.state])

  useEffect(() => {
    if (knownReposLoading.state.status !== LoadingStatus.Loaded) {
      return;
    }

    const params = match('/:repo/refs/:ref+/blob/:filepath+', location.state);
    console.debug("handling location state change", params);

    if (params.repo) {
      const repo = knownRepos.state.find((repo) => {
        return params.repo === repo.name;
      })

      console.debug("setting selectedRepo state", repo);
      selectedRepo.setState(repo);
    } else {
      selectedRepo.setState(undefined);
      return;
    }

    if (knownRefsLoading.state.status !== LoadingStatus.Loaded) {
      return;
    }

    let ref: Ref | undefined;

    if (params.ref) {
      ref ??= knownRefs.state.branches[params.ref];
      ref ??= knownRefs.state.tags[params.ref];

      console.debug("setting selectedRef state", ref);
      selectedRef.setState(ref);
    }

    if (ref) {
      console.debug("setting selectedRef state", ref);
      selectedRef.setState(ref);
    } else {
      const defaultBranchName =
        "main" in knownRefs.state.branches ? "main" :
          "master" in knownRefs.state.branches ? "master" :
            Object.keys(knownRefs.state.branches)[0] ?? Object.keys(knownRefs.state.tags)[0];

      location.replaceState(`/${params.repo}/refs/${defaultBranchName}`);
    }

    if (params.filepath) {
      selectedFilepath.setState(params.filepath);
    } else {
      selectedFileLoading.setIdle();
      selectedFile.setState(undefined);
      selectedFilepath.setState(undefined);
    }
  }, [location.state, knownRepos.state, knownRefs.state])

  return (
    <div className="flex-1 flex flex-col min-h-0 min-w-0" style={{
      fontFamily: "'JetBrains Mono', monospace",
      scrollbarWidth: "thin",
      scrollbarColor: "rgba(61,255,160,0.15) transparent",
    }}>
      <div className="flex flex-col h-dvh w-full overflow-hidden bg-background text-foreground">
        {knownReposLoading.state.status === LoadingStatus.Loading && (
          // <div className="absolute inset-0 flex flex-col items-center justify-center gap-2">
          //   <div className="w-4 h-4 border border-accent/30 border-t-accent rounded-full animate-spin" />
          //   <span className="text-[10px] text-muted-foreground">Loading…</span>
          // </div>

          <div className="flex items-center justify-center h-full gap-3">
            <div className="w-5 h-5 border border-accent/30 border-t-accent rounded-full animate-spin" />
            <span className="text-xs text-muted-foreground">
              Loading repositories…
            </span>
          </div>
        )}

        {knownReposLoading.state.status === LoadingStatus.Loaded && selectedRepo.state && (
          <Repo />
        )}

        {knownReposLoading.state.status === LoadingStatus.Error && (
          // <div className="absolute inset-0 flex flex-col items-center justify-center gap-2 px-4">
          //   <AlertTriangle size={16} className="text-destructive opacity-70" strokeWidth={1.5} />
          //   <span className="text-[10px] text-muted-foreground text-center leading-relaxed">
          //     {knownReposLoading.state.message}
          //   </span>
          // </div>

          <div className="flex flex-col items-center justify-center h-full gap-2">
            <AlertTriangle
              size={20}
              className="text-destructive opacity-70"
              strokeWidth={1.5}
            />
            <span className="text-xs text-destructive font-medium">
              Failed to load repositories
            </span>
            <span className="text-[11px] text-muted-foreground">
              {knownReposLoading.state.message}
            </span>
          </div>
        )}

        {knownReposLoading.state.status === LoadingStatus.Loaded && !selectedRepo.state && <>
          {/* ── Navbar ── */}
          <header className="shrink-0 flex items-center gap-0 border-b border-border bg-card">
            <Link className="flex items-center justify-center w-12 h-10 border-r border-border shrink-0" to="/">
              <Asterisk size={16} className="text-accent" strokeWidth={2} />
            </Link>
            <div className="relative flex items-center gap-2 h-10 border-r border-border">
              <div className="flex items-center gap-1.5 px-4 h-full ">
                <span className="text-xs text-muted-foreground">{window.location.host}</span>
              </div>
            </div>
          </header>
          <div className="max-w-3xl mx-auto px-8 py-10 h-dvh">
            <div className="flex items-baseline gap-3 mb-8">
              <h1 className="text-sm font-medium text-foreground">
                repos
              </h1>
              <span className="text-[10px] text-muted-foreground">
                {knownRepos.state.length} available
              </span>
            </div>

            <div className="divide-y divide-border border border-border">
              {knownRepos.state.map((repo) => (
                <Link key={repo.name} to={`/${repo.name}`} className="w-full flex items-center justify-between px-5 py-4 text-left hover:bg-secondary transition-colors group">
                  <div className="flex flex-col gap-1 min-w-0">
                    <span className="text-sm text-accent font-medium group-hover:text-accent/80 transition-colors truncate">
                      {repo.name}
                    </span>
                    {repo.description && (
                      <span className="text-xs text-muted-foreground truncate leading-relaxed">
                        {repo.description}
                      </span>
                    )}
                  </div>
                  <ArrowRight
                    size={13}
                    className="text-muted-foreground group-hover:text-accent shrink-0 ml-4 transition-colors"
                    strokeWidth={1.5}
                  />
                </Link>
              ))}
            </div>
          </div>
        </>}

        <footer className="shrink-0 flex items-center border-t border-border bg-card px-3 h-7 gap-4 text-[10px]">
          <div className="flex items-center gap-1.5">
            <Circle
              size={6}
              className={
                knownReposLoading.state.status === LoadingStatus.Error || selectedFileLoading.state.status === LoadingStatus.Error || currentTreeLoading.state.status === LoadingStatus.Error
                  ? "fill-destructive text-destructive"
                  : "fill-accent text-accent"
              }
            />
            <span className="text-muted-foreground">
              {knownReposLoading.state.status === LoadingStatus.Error || selectedFileLoading.state.status === LoadingStatus.Error || currentTreeLoading.state.status === LoadingStatus.Error
                ? "ERR"
                : selectedFileLoading.state.status === LoadingStatus.Loading || currentTreeLoading.state.status === LoadingStatus.Loading
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
                {selectedFileLoading.state.status === LoadingStatus.Loaded ? `${selectedFile.state!.timeElapsed.toFixed(1)}ms` : "—"}
              </span>
            </span>
          </div>

          <div className="flex-1" />

          <div className="flex items-center gap-3">
            {[ // TODO: Make these configurable
              { icon: Github, label: "GitHub", href: "https://github.com/iainjreid/source" },
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
                <Icon height={11} width={12} strokeWidth={1.5} />
                <span className="hidden sm:inline">{label}</span>
              </a>
            ))}
          </div>
        </footer>
      </div>
    </div>
  );
}

function convertTree(tree: Tree, selectedPath?: string) {
  const dirs: FileNode[] = Object.values(tree.dirs).map((entry) => ({
    kind: "dir",
    name: entry.name,
    path: entry.path,
    expanded: selectedPath?.startsWith(entry.path) ?? false,
    children: convertTree(entry.children, selectedPath),
  }));

  const files: FileNode[] = Object.values(tree.files).map((entry) => ({
    kind: "file",
    name: entry.name,
    path: entry.path,
  }));

  return [...dirs, ...files];
}
