import os
import sys
import yaml

def fix_sort_values_items(schema_dict):
    """
    If a schema has properties.sortValues (type: array), set items to allow
    string | integer | number via oneOf.
    Only overwrites when items are missing or currently 'object'.
    """
    if not isinstance(schema_dict, dict):
        return
    props = schema_dict.get("properties")
    if isinstance(props, dict) and "sortValues" in props:
        sv = props["sortValues"]
        if isinstance(sv, dict) and sv.get("type") == "array":
            items = sv.get("items")
            if not isinstance(items, dict) or items.get("type") == "object":
                sv["items"] = {
                    "oneOf": [
                        {"type": "string"},
                        {"type": "integer"},
                        {"type": "number"}
                    ]
                }

def main(input_file: str):
    with open(input_file, "r", encoding="utf-8") as f:
        spec = yaml.safe_load(f)

    # Walk components.schemas and patch sortValues where present
    components = spec.get("components", {})
    schemas = components.get("schemas", {})
    for _, schema in schemas.items():
        fix_sort_values_items(schema)

    # write with suffix
    base, ext = os.path.splitext(input_file)
    output_file = f"{base}-fix-sortvals{ext}"
    with open(output_file, "w", encoding="utf-8") as f:
        yaml.dump(spec, f, sort_keys=False, allow_unicode=True)

    print(f"Wrote: {output_file}")

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python 3_fix_sortvalues.py <openapi.yaml>")
        sys.exit(1)
    main(sys.argv[1])