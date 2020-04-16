# bump-card

Convenience tool to put a timestamp in frontmatter of markdown files.

This can be used to keep documentation up to date, where you can sort by deprecated files. \
When having reviewed a document, being satisfied that it is up to date, you can then bump it.

## Usage

```
Usage: ./bump-card [options] [files...]

Options:
   -h --help       Print help
   -b --bump       Bump files
   -r --recursive  List subdirectories recursively
   -f --force      Bump even if no previous frontmatter

If no files given, list all files
```
