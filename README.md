# LibGen OPDS Bridge

## Description
This creates a bridge between Library Genesis and OPDS. You can search and download books on LibGen via KOReaders OPDS Search. Most read lists from GoodReads are also parsed, and allow you to directly jump to a relevant LibGen search.

## Limitations
- Only 25 results / search

## KUAL Extension Installation
- Get the [current release here](https://gitea.va.reichard.io/evan/libgen-opds/releases) (e.g. `libgen-opds_0_0_3.zip`)
- Extract and add the `libgen-opds` folder to the `extensions` folder on the Kindle.
- Open KUAL -> Start LibGen OPDS Bridge
- Open KOReader and add a new OPDS Catalog in KOReader: Search -> OPDS catalog -> "+" (Upper Left):
  - Catalog Name: LibGen Fiction
  - Catalog URL: http://127.0.0.1:5144
  
## Building
The output will be in `./build/dist`:

    # Kindle KUAL Extension
    make build_kual_extension

    # MultiArch Binaries
    make build_multiarch

    # Docker MultiArch (Note: You may need to use `docker buildx create --use`)
    make build_multiarch_docker

    # Docker
    make build_docker

## To Do
- [ ] Better search results descriptions (file size, upload date, ?)
- [ ] Configuration (port, etc)
- [ ] Docker support
- [ ] Pin DNS server ([example](https://koraygocmen.medium.com/custom-dns-resolver-for-the-default-http-client-in-go-a1420db38a5d))
