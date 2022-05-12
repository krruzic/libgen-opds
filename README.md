# LibGen OPDS Bridge

## Description
This creates a bridge between Library Genesis and OPDS. You can search and download books on LibGen via KOReaders OPDS Search. Most read lists from GoodReads are also parsed, and allow you to directly jump to a relevant LibGen search.

## Limitations
At the moment, only `epub` books are searched for, results are limited to 25 items, and only the Fiction category is searched.

## Installation
- Add the `libgen-opds` folder in `extensions` to the `extensions` folder on the Kindle.
- Open KUAL -> Start LibGen OPDS Bridge
- Open KOReader and add a new OPDS Catalog in KOReader: Search -> OPDS catalog -> "+" (Upper Left):
  - Catalog Name: LibGen Fiction
  - Catalog URL: http://127.0.0.1:5144
  
## Setup Build Environment

    docker-compose build

## Build For Kindle

    docker-compose run --rm libgen-opds compile_kindle
