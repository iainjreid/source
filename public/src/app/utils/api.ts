export interface Repos {
  repos: Repo[];
}

export interface Repo {
  id: string;
  name: string;
  description: string;
}

export interface ApiRefs {
  branches: Record<string, string>;
  tags: Record<string, string>;
}

export interface Refs {
  branches: Record<string, Ref>;
  tags: Record<string, Ref>;
}

export enum RefKind {
  Branch,
  Tag,
}

export interface Ref {
  name: string;
  shortName: string;
  hash: string;
  kind: RefKind
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

export function GetRepos() {
  console.debug("loading repos");

  return once<Record<string, Repo>>('GetRepos', () => {
    return request<Repos>('/repos').then((res) => {
      return res.repos.reduce((acc, repo) => ({
        [repo.name]: repo,
        ...acc,
      }), {});
    });
  });
}

export function GetRefs(repo: Repo): Promise<ApiRefs> {
  console.debug("loading refs", repo);

  return once(`GetRefs#${repo.name}`, () => {
    return request<ApiRefs>(`/${repo.id}/refs`).then(tap((data) => {
      console.log("Refs loaded", data);
    }));
  });
}

export function GetTree(repo: Repo, ref: Ref): Promise<Tree> {
  console.debug("loading tree", repo, ref);

  return once(`GetTree#${repo.name}/${ref.name}`, () => {
    return request<Tree>(`/${repo.id}/tree/${ref.hash}`).then(tap((data) => {
      console.log("Tree loaded", data);
    }));
  });
}

export function GetFile(repo: Repo, ref: Ref, path: string): Promise<File> {
  console.debug("loading file", repo, ref, path);

  return once(`GetRefs#${repo.name}/${ref.name}/${path}`, () => {
    return request<File>(`/${repo.id}/blob/${ref.hash}/${path}`).then(tap((data) => {
      console.log("File loaded", data);
    }));
  });
}

function request<T>(url: string): Promise<T> {
  const controller = new AbortController();
  const timeout = setTimeout(() => {
    controller.abort("Request timed out")
  }, 15_000);

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

  return promise;
}

const cache = new Map<string, Promise<any>>();

function once<T>(key: string, fn: () => Promise<T>): Promise<T> {
  if (!cache.has(key)) {
    cache.set(key, fn());
  }

  return cache.get(key)!;
}

function tap<T>(fn: (data: T) => any): (data: T) => T {
  return (data: T) => {
    fn(data);
    return data
  };
}
