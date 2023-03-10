# Linkchecker

Checks links in a powerpoint file for reachability

## Usage

```bash
linkchecker check filename.pptx
```

## Example

```bash
 go run main.go check testdata/someuris.pptx
❌ Target https://www.john-doe.com/ is not reachable in slide: slide3
☑️ Target https://www.google.de/ is reachable in slide: slide3
❌ Target https://www.john-doe.com/ is not reachable in slide: slide2
☑️ Target https://www.google.de/ is reachable in slide: slide2
```