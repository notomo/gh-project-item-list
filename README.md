gh extension to filter project items by custom field.
This will not be necessary if the official api allows project's filter string.

## Usage

```
gh project-item-list -project-url=https://github.com/users/notomo/projects/1 -limit=20 -jq='.[] | select(.fieldValues.nodes|any(.field.name == "Status" and .name == "Todo"))'
```
