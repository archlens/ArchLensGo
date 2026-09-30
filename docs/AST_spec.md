# AST specification for parsers
ArchLens is a Go based dependency visualizer. Finding dependencies for the individual languages are handled by language specific parsers. These are called on individual files and are used to build a graph which then gets converted to a mermaid diagram.

## Input
The parser takes a filename and a potential flag for root directory of the project which defaults to `.`.

The file being parsed can be at any depth below the root. Do not assume it sits in the root directory.

## Specification
Parsers should read a single file. Generate an AST based on the imports found in that file. Compile the AST and filter imports to only internal dependencies, then package this info into a json string with the format below. Send data over stdout. DO NOT SORT THE DEPENDENCIES THERE IS NO POINT IN DOING THIS.

```json
{
  "package": string,
  "dependencies": [
    string,
    string
  ]
}
```

### Resolving imports
An import is internal if and only if it resolves to a file or module located inside the project root. To decide this, resolve every import the way the language's own toolchain would, not by checking against the root alone:

- Check every location the language searches for imports. This always includes the project root and the directory of the file being parsed, plus any other rule the language defines (relative imports, parent directories, source directories, and so on).
- Handle relative imports by resolving them against the location of the file being parsed.
- An import that resolves to a file inside the root is internal, no matter which search location found it. An import that resolves outside the root, or does not resolve at all, is external and must be left out (standard library, third-party packages, unresolvable names).
- If the language allows importing a name that could be either a module or a member of a module, check whether it resolves to a module first and fall back to the enclosing module.


### Naming
Dependencies are in the format `package.subpackage.filename`, where `subpackage` is only present if a subpackage is specified. Use the language's usual separator.

- Always build the name from the resolved file's path relative to the project root, never from the text of the import statement. The same file must produce the same string whether it was imported by an absolute, relative, or sibling-style import.
- Build `package` for the parsed file the same way, from the location of the file relative to the root. A file directly in the root has an empty package.
- These names must line up with the `package` value the parser reports when it is run on the dependency itself, so graph edges connect correctly.

### Output rules
- Do not include duplicates. Keep dependencies in the order they are found, without sorting.
- Do not list the parsed file as a dependency of itself.
- Only JSON goes to stdout. Errors (unreadable file, syntax error, file outside the root) go to stderr with a non-zero exit code.

### Cases to verify
A parser is not done until it handles each of these:

1. A file in the root importing another file in the root.
2. A file in a subdirectory importing a file in that same subdirectory, with the root left at the default.
3. A file importing something from a parent or sibling directory using relative syntax, if the language has it.
4. An import of a whole package or directory, not just a single file.
5. A mix of internal and external imports, where only the internal ones are returned.
6. The same internal file imported in two different ways, which appears once in the output.