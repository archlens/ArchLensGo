import ast
import json
import sys

def node_to_dict(node):
    if not isinstance(node, ast.AST):
        return node
    result = {"type": type(node).__name__}
    if hasattr(node, "lineno"):
        result["line"] = node.lineno
    children = []
    for field, value in ast.iter_fields(node):
        if isinstance(value, list):
            for item in value:
                if isinstance(item, ast.AST):
                    children.append(node_to_dict(item))
        elif isinstance(value, ast.AST):
            children.append(node_to_dict(value))
        elif value is not None:
            result.setdefault("value", str(value))
    if children:
        result["children"] = children
    return result

if __name__ == "__main__":
    path = sys.argv[1]
    with open(path) as f:
        tree = ast.parse(f.read(), filename=path)
    print(json.dumps(node_to_dict(tree)))