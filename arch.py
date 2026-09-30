#!/usr/bin/env python3
"""ArchLens parser for Python.

Reads a single .py file, builds an AST, extracts imports, keeps only those that
resolve to a file or directory inside the project root, and prints JSON to
stdout:

    {"package": str, "dependencies": [str, ...]}

Resolution follows Python's own lookup rules:
  * absolute imports are searched in the directory of the parsed file (Python
    puts a script's own directory on sys.path), the project root, and
    <root>/src if it exists;
  * relative imports are resolved against the location of the parsed file;
  * regular modules/packages win over bare directories (namespace packages),
    and bare directories never shadow a standard library module.

Names are always built from the resolved path relative to the root, e.g.
`app.core.models`, so the same file gets the same string however it was
imported, and matches the `package` value of that file's own parser run.

Usage:
    python archlens_parser.py path/to/file.py [-r path/to/project]
"""
from __future__ import annotations

import argparse
import ast
import json
import sys
from pathlib import Path

STDLIB = getattr(sys, "stdlib_module_names", frozenset())


def dotted(target: Path, root: Path) -> str | None:
    """Dotted name of a file or directory relative to root (None if outside)."""
    try:
        parts = list(target.relative_to(root).parts)
    except ValueError:
        return None
    if parts and parts[-1].endswith(".py"):
        parts[-1] = parts[-1][:-3]
    if parts and parts[-1] == "__init__":
        parts.pop()
    return ".".join(parts) or None


def locate(module: str, dirs: list[Path], root: Path) -> str | None:
    """Resolve a dotted module against search dirs; return its root-relative name."""
    parts = module.split(".") if module else []

    # 1. Regular modules and packages, in search order.
    for base in dirs:
        p = base.joinpath(*parts)
        candidates = [p / "__init__.py"]
        if parts:
            candidates.insert(0, p.parent / (p.name + ".py"))
        for c in candidates:
            if c.is_file():
                name = dotted(c, root)
                if name:
                    return name

    # 2. Bare directories (namespace packages), never shadowing the stdlib.
    if parts and parts[0] not in STDLIB:
        for base in dirs:
            p = base.joinpath(*parts)
            if p.is_dir():
                name = dotted(p, root)
                if name:
                    return name
    return None


def longest_internal_prefix(module: str, dirs: list[Path], root: Path) -> str | None:
    """`import a.b.c` -> a.b.c if internal, else the longest internal prefix."""
    parts = module.split(".")
    for i in range(len(parts), 0, -1):
        found = locate(".".join(parts[:i]), dirs, root)
        if found:
            return found
    return None


def within(path: Path, root: Path) -> bool:
    return path == root or root in path.parents


class ImportCollector(ast.NodeVisitor):
    """Collects internal dependencies in source order, without duplicates."""

    def __init__(self, root: Path, file_dir: Path):
        self.root = root
        self.file_dir = file_dir
        dirs = [file_dir, root]
        if (root / "src").is_dir():
            dirs.append(root / "src")
        self.abs_dirs = list(dict.fromkeys(dirs))  # dedupe, keep order
        self.deps: dict[str, None] = {}  # insertion-ordered set

    def add(self, name: str | None) -> None:
        if name:
            self.deps[name] = None

    def visit_Import(self, node: ast.Import) -> None:
        for alias in node.names:
            self.add(longest_internal_prefix(alias.name, self.abs_dirs, self.root))

    def visit_ImportFrom(self, node: ast.ImportFrom) -> None:
        if node.level == 0:
            if not node.module:
                return
            base, dirs = node.module, self.abs_dirs
        else:
            base_dir = self.file_dir
            for _ in range(node.level - 1):
                base_dir = base_dir.parent
            if not within(base_dir, self.root):
                return  # relative import climbs out of the project
            base, dirs = node.module or "", [base_dir]

        base_hit = locate(base, dirs, self.root)
        for alias in node.names:
            # `from pkg import name`: module first, fall back to the enclosing module
            sub = None
            if alias.name != "*":
                sub = locate(f"{base}.{alias.name}" if base else alias.name, dirs, self.root)
            self.add(sub or base_hit)


def main() -> int:
    ap = argparse.ArgumentParser(description="ArchLens Python parser")
    ap.add_argument("file", type=Path)
    ap.add_argument("-r", "--root", type=Path, default=Path("."),
                    help="project root used to decide what is internal (default: .)")
    args = ap.parse_args()

    root = args.root.resolve()
    path = args.file.resolve()
    if not root.is_dir():
        print(f"error: root {root} is not a directory", file=sys.stderr)
        return 1
    if not within(path, root):
        print(f"error: {path} is not under root {root}", file=sys.stderr)
        return 1

    try:
        tree = ast.parse(path.read_bytes(), filename=str(path))
    except (SyntaxError, ValueError, OSError) as e:
        print(f"error: cannot parse {path}: {e}", file=sys.stderr)
        return 1

    package = ".".join(path.relative_to(root).parent.parts)
    self_module = dotted(path, root)

    collector = ImportCollector(root, path.parent)
    collector.visit(tree)
    deps = [d for d in collector.deps if d != self_module]

    json.dump({"package": package, "dependencies": deps}, sys.stdout)
    sys.stdout.write("\n")
    return 0


if __name__ == "__main__":
    sys.exit(main())