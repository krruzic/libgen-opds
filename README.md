# LibGen OPDS Bridge

## Description
This creates a bridge between Library Genesis and OPDS. You can search and download books on LibGen via KOReaders OPDS Search. Most read lists from GoodReads are also parsed, and allow you to directly jump to a relevant LibGen search.

## Limitations
- Only `epub`
- Only fiction
- Only 25 results / search

## KUAL Extension Installation
- Get the [current release here](https://gitea.va.reichard.io/evan/libgen-opds/releases) (e.g. `libgen-opds_0_0_2.zip`)
- Extract and add the `libgen-opds` folder to the `extensions` folder on the Kindle.
- Open KUAL -> Start LibGen OPDS Bridge
- Open KOReader and add a new OPDS Catalog in KOReader: Search -> OPDS catalog -> "+" (Upper Left):
  - Catalog Name: LibGen Fiction
  - Catalog URL: http://127.0.0.1:5144
  
## Building

    # Kindle KUAL Extension
    make build_kual_extension

    # Docker MultiArch (Note: You may need to use `docker buildx create --use`)
    make build_docker

## To Do
- [ ] Better search results descriptions (file size, upload date, ?)
- [ ] Configuration (port, etc)
- [ ] Docker support
- [ ] Logging (any at all...)
- [ ] More sources (scimag, etc)
- [ ] Multi format support
- [ ] Pin DNS server ([example](https://koraygocmen.medium.com/custom-dns-resolver-for-the-default-http-client-in-go-a1420db38a5d))
