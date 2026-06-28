import { useEffect, useRef } from "react";

declare const Prism: typeof import("prismjs")

const extensionMap: Record<string, string> = {
  go: "go",
  js: "javascript",
  mjs: "javascript",
  cjs: "javascript",
  jsx: "jsx",
  ts: "typescript",
  tsx: "tsx",
  json: "json",
  html: "markup",
  xml: "markup",
  svg: "markup",
  css: "css",
  scss: "scss",
  sh: "bash",
  bash: "bash",
  zsh: "bash",
  py: "python",
  rb: "ruby",
  rs: "rust",
  java: "java",
  php: "php",
  yaml: "yaml",
  yml: "yaml",
  toml: "toml",
  sql: "sql",
  md: "markdown",
  c: "c",
  h: "c",
  cpp: "cpp",
  cc: "cpp",
  cxx: "cpp",
  hpp: "cpp",
};

export function resolveLanguage(filePath: string) {
  const basename = filePath.split("/").pop()!.toLowerCase();

  switch (basename) {
    case "dockerfile":
      return "docker";

    case "makefile":
    case "gnumakefile":
      return "makefile";

    case ".gitignore":
      return "gitignore";

    case ".editorconfig":
      return "editorconfig";
  }

  const ext = basename.split(".").pop()!;

  return extensionMap[ext] ?? "plain";
}

export function Code({ filePath, lines }: { filePath: string, lines: string[] }) {
  const codeRef = useRef(null);
  const dataPad = "0".repeat(Math.ceil(Math.log10(lines.length + 1)));

  const language = resolveLanguage(filePath);

  useEffect(() => {
    if (codeRef.current) {
      Prism.highlightElement(codeRef.current);
    }
  }, [lines, language]);

  return (
    <div className={`flex min-h-full w-full`}>
      <div className="line-numbers select-none text-right px-4 text-gray-500">
        {lines.map((_, i) => (
          <span key={i} data-pad={dataPad}></span>
        ))}
      </div>
      <pre className={`pl-4 flex min-h-full w-full`}>
        <code ref={codeRef} className={`min-h-full w-full language-${language}`}>
          {lines.join("\n")}
        </code>
      </pre>
    </div>
  );
}
