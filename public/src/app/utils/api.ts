
export interface Refs {
  branches: Record<string, string>;
  tags: Record<string, string>;
}

export interface Tree {
  dirs: Record<string, TreeEntry>;
  files: Record<string, TreeEntry>;
}

export interface TreeEntry {
  name: string;
  path: string;
  isFile: boolean;
  children: Tree;
}

export interface File {
  fileName: string;
  fileBytes: number;
  fileLines: string[];
  lineCount: number;
  lastCommitHash: string;
  lastCommitMsg: string;
  lastCommitAuthor: string;
  timeElapsed: number;
}

export function GetFile(repoId: string, ref: string, path: string): APIRes<File> {
  return request<File>(`http://localhost:8080/${repoId}/blob/${ref}/${path}`);
}

interface APIRes<T> {
  destructor: () => void
  promise: Promise<T>
}

function request<T>(url: string): APIRes<T> {
  const controller = new AbortController();
  const timeout = setTimeout(() => {
    controller.abort("Request timed out")
  }, 15_000);

  const destructor = () => {
    clearTimeout(timeout)
    controller.abort()
  }

  const promise = fetch(url, { signal: controller.signal })
    .then((res) => {
      if (!res.ok) {
        return Promise.reject(`HTTP ${res.status} — ${res.statusText}`);
      }

      return res.json()
    })
    .finally(() => {
      clearTimeout(timeout);
    });

  return {
    destructor,
    promise,
  }
}
