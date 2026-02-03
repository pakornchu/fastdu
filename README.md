# fastdu

Quickly get file sizes from pattern file.

## Build
```
go build -o fastdu fastdu.go
```

## Usage
```
fastdu [-r|--raw] [-s|--summary] [-w|--worker N] PATTERNFILE
  -r, --raw          Show raw bytes
  -s, --summary      Show summary only
  -v, --version      Show version
  -w, --worker int   Number of worker (default 4)
```

## Notes
- Files that cannot be stat will be ignore without any error output
