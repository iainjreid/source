import { Dispatch, SetStateAction, useEffect, useState } from "react";
import { Ref, Refs, Repo } from "../utils/api";

export type FileResponse = {
  fileName: string;
  fileBytes: number;
  fileLines: string[];
  lineCount: number;
  lastCommitHash: string;
  lastCommitMsg: string;
  lastCommitAuthor: string;
  timeElapsed: number;
};

export enum LoadingStatus {
  Idle,
  Loading,
  Loaded,
  Error,
}

type LoadingState =
  | { status: LoadingStatus.Idle }
  | { status: LoadingStatus.Loading }
  | { status: LoadingStatus.Loaded }
  | { status: LoadingStatus.Error; message: string };

export type FileNode =
  | { kind: "dir"; name: string; path: string; children: FileNode[]; expanded: boolean }
  | { kind: "file"; name: string; path: string };

export const useKnownRepos = createSharedState<Repo[]>([]);
export const useKnownReposLoading = createSharedLoadingState();

export const useSelectedRepo = createSharedState<Repo | undefined>(undefined);

export const useKnownRefs = createSharedState<Refs>({
  branches: {},
  tags: {},
});

export const useKnownRefsLoading = createSharedLoadingState();

export const useSelectedRef = createSharedState<Ref | undefined>(undefined);

export const useCurrentTree = createSharedState<FileNode[]>([]);
export const useCurrentTreeLoading = createSharedLoadingState();

export const useSelectedFilepath = createSharedState<string | undefined>(undefined);

export const useSelectedFile = createSharedState<FileResponse | undefined>(undefined);
export const useSelectedFileLoading = createSharedLoadingState();

export interface SharedState<T> {
  state: T;
  setState: Dispatch<SetStateAction<T>>;
}

export function createSharedState<T>(initial: T): () => SharedState<T> {
  const store = {
    value: initial,
    events: new EventTarget(),
  };

  return (): SharedState<T> => {
    const [state, setInternal] = useState(store.value);

    useEffect(() => {
      console.debug("propagating shared state");
      const update = () => setInternal(store.value);

      store.events.addEventListener("change", update);

      return () => {
        store.events.removeEventListener("change", update);
      };
    }, []);

    const setState = (next: T | ((prev: T) => T)) => {
      store.value =
        typeof next === "function"
          ? (next as (prev: T) => T)(store.value)
          : next;

      store.events.dispatchEvent(new Event("change"));
    };

    return { state, setState };
  }
}

export interface Loadable {
  state: LoadingState;
  setIdle: () => void
  setLoading: () => void
  setLoaded: () => void
  setError: (message: string) => void
}

export function createSharedLoadingState(): () => Loadable {
  const useState = createSharedState<LoadingState>({ status: LoadingStatus.Idle });

  return (): Loadable => {
    const { state, setState } = useState();

    function setIdle() {
      setState({ status: LoadingStatus.Idle })
    }

    function setLoading() {
      setState({ status: LoadingStatus.Loading })
    }

    function setLoaded() {
      setState({ status: LoadingStatus.Loaded })
    }

    function setError(message: string) {
      setState({ status: LoadingStatus.Error, message })
    }

    return { state, setIdle, setLoading, setLoaded, setError };
  }
}
