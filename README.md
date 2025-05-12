# Usage

1. Build the project using `make build`
2. Run the program using the command:
```bash
bin/synchro sync --source_dir <src> --destination_dir <dst>
```
shorthand flags are also supported:
```bash
bin/synchro sync --src <src> --dst <dst>
```

`--delete_missing` optional flag can be used to delete files in the destination directory that are not present in the source directory anymore
long version:
```bash
bin/synchro sync --source_dir <src> --destination_dir <dst> --delete_missing
```
shorthand version:
```bash
bin/synchro sync --src <src> --dst <dst> -d
```

