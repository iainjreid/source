import { FC } from "react"
import { ChevronDown, ChevronRight, File, Folder, FolderOpen } from "lucide-react";

interface TreeNodeProps {
  node: FileNode;
  depth: number;
  selectedPath: string;
  onSelect: (path: string) => void;
  onToggle: (path: string) => void;
}

export const TreeNode: FC<TreeNodeProps> = (props) => {
  const indent = props.depth * 14;

  if (props.node.kind === "file") {
    const isSelected = props.selectedPath === props.node.path;
    return (
      <button
        onClick={() => props.onSelect(props.node.path)}
        className="w-full flex items-center gap-1.5 py-[3px] pr-3 text-left transition-colors group"
        style={{ paddingLeft: `${indent + 8}px` }}
      >
        <File size={12} className={isSelected ? "text-accent" : "text-muted-foreground"} strokeWidth={1.5} />
        <span className={`text-xs font-mono truncate transition-colors ${isSelected ? "text-accent font-medium" : "text-muted-foreground group-hover:text-foreground"
          }`}>
          {props.node.name}
        </span>
      </button>
    );
  }

  return (
    <div>
      <button
        onClick={() => props.onToggle(props.node.path)}
        className="w-full flex items-center gap-1.5 py-[3px] pr-3 text-left hover:text-foreground transition-colors group"
        style={{ paddingLeft: `${indent + 8}px` }}
      >
        {props.node.expanded
          ? <ChevronDown size={11} className="text-muted-foreground shrink-0" strokeWidth={2} />
          : <ChevronRight size={11} className="text-muted-foreground shrink-0" strokeWidth={2} />
        }
        {props.node.expanded
          ? <FolderOpen size={12} className="text-accent/60 shrink-0" strokeWidth={1.5} />
          : <Folder size={12} className="text-muted-foreground shrink-0" strokeWidth={1.5} />
        }
        <span className="text-xs font-mono text-foreground/80 group-hover:text-foreground truncate transition-colors">
          {props.node.name}
        </span>
      </button>
      {props.node.expanded && (
        <div>
          {props.node.children.map((child) => (
            <TreeNode
              key={child.path}
              node={child}
              depth={props.depth + 1}
              selectedPath={props.selectedPath}
              onSelect={props.onSelect}
              onToggle={props.onToggle}
            />
          ))}
        </div>
      )}
    </div>
  );
}
