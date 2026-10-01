# AST specification for parsers

ArchLens is a Go based dependency visualizer. Finding dependencies for the individual languages is handled by language specific parsers. A parser is called on one file at a time. ArchLens combines the results into a graph whose nodes are packages and whose edges mean "a file in package A imports package B", then converts the graph to a mermaid diagram.

The specification is language independent. Every parser, whatever the language, produces the same output format and the same names, so results from different parsers can be combined.

## Terms

- **Project root**: the directory the parser is told to treat as the top of the project. Only things inside it can be internal.
- **Source file**: a file the parser understands (for example `.py` for the Python parser).
- **Package**: the directory containing a source file. This is the node ArchLens draws, regardless of whether the language itself calls the directory a package, module or something else.
- **Internal**: an import that resolves to a file or directory inside the project root and is part of the project's own code (see "Resolving imports").

## Input

```console
<parser> <file> [-r DIR | --root DIR]
```

- `file` is a path to one source file. It is resolved against the **working directory**, not against the root.
- `--root` is the project root and defaults to `.` (the working directory). Resolve it to an absolute path before using it.
- Resolve both paths through symlinks before comparing them. The file must be inside the root, at any depth. Do not assume it sits directly in the root.
- ArchLens runs parsers with the working directory set to the root, an absolute `--root`, and the file path relative to the root. A parser must still behave correctly when run any other way.

## Output

Print exactly one JSON object followed by a newline to stdout, UTF-8 encoded:

```json
{
  "package": "app/core",
  "dependencies": ["app", "app/utils"]
}
```

- `package` is the package of the parsed file.
- `dependencies` is the list of internal packages the file imports. Always present. When there are none, it is `[]`, never `null` and never omitted.
- There are no other fields.

## What counts as an import

- Every **static** import in the file counts, wherever it appears: at the top, inside a function, inside a conditional or `try` block, in a type-checking-only block.
- **Dynamic** imports, where the module name is computed at runtime (`importlib.import_module(name)`, `require(variable)` and similar), are ignored.
- Imports of the language's standard library and of third-party packages are external and left out (see below).

## Resolving imports

An import is internal if and only if it resolves to a file or directory inside the project root *and* is not a third-party install. To decide this, resolve every import the way the language's own toolchain would, not by checking against the root alone:

- Check every location the language searches for imports. This always includes the project root and the directory of the file being parsed, plus any other rule the language defines (relative imports, parent directories, source directories, and so on).
- Handle relative imports by resolving them against the location of the file being parsed. A relative import that climbs out of the root is external.
- An import that resolves inside the root is internal, no matter which search location found it. An import that resolves outside the root, or does not resolve at all, is external and must be left out (standard library, third-party packages, unresolvable names).
- **Third-party installs inside the root are external.** Directories that a package manager fills with other people's code (a virtual environment's `site-packages`, `node_modules`, `vendor` directories, and similar) are inside the root but are not the project's own code. Imports that resolve into them are external. Each parser documents which directories it treats this way.
- If the language allows importing a name that could be either a module or a member of a module, check whether it resolves to a module first and fall back to the enclosing module.
- If importing a package only re-exports names from other files, the dependency is the package that was imported, not the package where the name was originally defined.
- Decide "inside the root" **after** resolving symlinks. A symlink inside the root that points outside it is external. Several paths that lead to the same real file give the same result.

### Naming

A *package* is the directory containing a source file, expressed as its path
relative to the project root with `/` as the separator on every platform.
The root directory itself is `.`.

- `package` is the package of the parsed file.
- Each entry in `dependencies` is the package of the resolved file that was
  imported (never the text of the import statement, and never a filename).
- A dependency that resolves to the parsed file's own package is omitted.
- The same string must be produced for a package whether it is the parsed
  file's own package or a dependency of another file.

## Output rules

- No duplicates. Because ArchLens counts one import per file between two packages, a file that imports package B five times contributes exactly one `B`. The weight of an edge in the diagram is the number of files in package A that import package B.
- Keep dependencies in the order they are first found. Do not sort them. The order does not change the diagram, but it must be deterministic so test expectations are stable and diffs stay readable.
- Only the JSON object goes to stdout, and nothing is written to stdout on failure.

## Errors

Errors go to stderr as a short message, with a non-zero exit code and nothing on stdout. This covers an unreadable file, a syntax error, a file outside the root, and a root that is not a directory. An import that cannot be resolved is not an error: it is external and is left out.

## Cases to verify

A parser is not done until it handles each of these. Keep them as fixtures (see below).

**Basic resolution**

1. A file in the root importing another file in the root. The result is empty, since the dependency is the file's own package (`.`).
2. A file in a subdirectory importing a file in that same subdirectory, with the root left at the default. Also empty.
3. A file importing something from a parent or sibling directory using relative syntax, if the language has it.
4. An import of a whole package or directory, not just a single file.
5. A mix of internal and external imports, where only the internal ones are returned.
6. The same package imported in two or more different ways (absolute, relative, whole package, single file), which appears once in the output.
**Edges and boundaries**

7. A file with no imports: `"dependencies": []`.
8. A file in the root importing a package in a subdirectory, and a subdirectory file importing something in the root (package `.` as a dependency).
9. Two files that import each other, in different packages (a cycle). Each reports the other's package.
10. An unresolvable import mixed with valid ones. Only the valid internal ones are returned.
11. An import that resolves through a symlink to a file outside the root: external.
12. An import that resolves into an installed third-party directory inside the root (`node_modules`, a virtual environment): external.
13. A relative import that climbs out of the root: external.
14. An import inside a function or a type-checking-only block: counted. A dynamic import: ignored.
15. A package whose entry file re-exports from submodules: the dependency is the imported package.
**Failures**

16. A syntax error: error on stderr, non-zero exit, nothing on stdout.
17. A file outside the root: error on stderr, non-zero exit, nothing on stdout.

## Conformance fixtures

Keep one fixture directory per case under `testdata/<language>/<case>/`, containing a small project and an `expected.json` for each file the case checks. A test harness in the ArchLens repository runs the language's parser on each file and compares its output to the expected JSON, byte for byte after JSON parsing. A new parser is correct when it passes every fixture, so it can be written or generated from this document alone.
