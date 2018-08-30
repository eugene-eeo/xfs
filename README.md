# <img src='static/leaf.png' width='30px'/> xfs

xfs is a set of tools to enable full text search on various types of files.
On it's own it ships with a very minimal set of binaries that are meant to
be composed to be the core of something like macOS's Spotlight (but with less
features). Currently supports macOS and _probably_ Linux.

1. [Installation](#installation)
1. [Quickstart](#quickstart)
1. [Configuration Options](#configuration-options)
1. [Snippets](#snippets)


## Installation

xfs requires the [Go compiler](https://golang.org/). After that you can do:

```
$ go get github.com/eugene-eeo/xfs/...
```

This installs a few binaries including:

 - `xfs-watch` – watches over directories and emits events to stdout.
 - `xfs-dispatch` – takes in events from stdin and handles them.
 - `xfs-index` – a command line tool to interact with the search index.
 - `xfs-search` – returns a list of matching files given some query.
 - `xfs-pdf` - extracts text from PDFs and dumps them to stdout.


## Quickstart

### Initial config

First create `~/.xfsrc`. It is a configuration file that is written in a
friendlier version of JSON which allows for JS `//` and `/* ... */` style
comments, as well as trailing commas.

```js
{
  "watch": [
    "~/Documents/...",
    "~/notes/",
  ],
  // comments are allowed
  "dispatch": [
    ["application/pdf",    "xfs-pdf"],
    ["application/x-text", "cat"],
  ],
}
```

This is a basic setup that allows for `~/Documents` to be recursively
watched, e.g. if you modified `~/Documents/foo/bar.pdf` then it will
be handled, and also for `~/notes` to be watched non recursively (i.e.
if `~/notes/foo/bar` was modified then no event is emitted).

### Starting up the watchers

Run the following in a different terminal:

```sh
$ xfs-watch | xfs-dispatch
```

This starts up the `xfs-dispatch` binary. At any point, you can hit
**Ctrl-C** when you feel like you've had enough. You should find that
when you edit files in either `~/Documents` or `~/notes` an event will
be emitted. To force an index of the watched directories you can do:

```sh
$ cd ~/Documents && find . -exec touch {} +
$ cd ~/notes     && find . -exec touch {} +
```

### Actually Searching

So far we've done a lot of scaffolding but very little actual searching.
If you want to search for avocados for example, you can simply do:

```sh
$ xfs-search 'avocados'
```

Note that the filenames themselves are also indexed, so you can search by
filename equally as easily – should you have an `avocado.txt` for example
you could also find that. If you want to do more advanced queries you can
refer to the [Bleve Query String documentation](http://blevesearch.com/docs/Query-String-Query/).

### Handling more file formats

Currently with our config, we're only indexing files which are PDFs or text
files. We can do much more than that, if you install the right tools. For
example say you've installed a tool, say `doc2text` that converts Microsoft
Word documents into text. You could add an entry into the `dispatch` array:

```js
{
  "dispatch": [
    ["application/msword", "doc2text"],
    ["application/pdf",    "xfs-pdf"],
    ["application/x-text", "cat"],
  ],
}
```

This means that when for example `~/notes/foo.doc` changes, `xfs-dispatch`
will now index the output of `doc2text ~/notes/foo.doc` – i.e. `xfs-dispatch`
does the equivalent of `doc2text ~/notes/foo.doc | xfs-index set ~/notes/foo.doc`.
The `dispatch` supports globs as well, so you could this for example:

```js
{
  "dispatch": [
    ["application/font-*", "myfonttool"]
    ["application/msword", "doc2text"],
    ["application/pdf",    "xfs-pdf"],
    ["application/x-text", "cat"],
  ]
}
```

Now `myfonttool` will be used to index TTF, WOFF, WOFF2, or OTF files. To
see the full list of supported filetypes you can visit this [README](https://github.com/h2non/filetype).
Additionally, two more mimetypes are added: **application/x-text** which is
used when the file is not one of the listed mimetypes and does not contain
binary content, and **unknown** when the file contains binary content.
Beware however that the dispatching algorithm is essentially this:

```
foreach [glob, handler] in the dispatch array:
    if glob matches mimetype:
        index file using handler
        break
```

So you could run into problems if you match against something too general
before other more specific tools. All paths are passed in full (i.e. they
are all absolute paths) to the handlers.


## Configuration Options

| Key            | Type               | Default | Meaning                        |
|:--------------:|--------------------|:-------:|--------------------------------|
| **watch**      | `array<string>`    | `[]`    | Array of watched directories. Directories that end in `/...` are interpreted similarly to how `go get` works and are indexed recursively. |
| **ignore**     | `array<string>`    | `[]`    | Array of ignored directories. They are recursive, for instance if you ignore `~/notes/pdfs`, then a change in `~/notes/pdfs/subdir/a.txt` will be handled. |
| **poll**       | `int`              | `1`     | Polling interval for filesystem changes (in seconds) |
| **data_dir**   | `string`           | `~/.xfs`| Directory where the index file and other information are stored. |
| **dispatch**   | `array<[2]string>` | `[]`    | Array of _entries_, which are arrays of size 2 containing a [glob](https://github.com/gobwas/glob) and it's corresponding handler. When dispatching the first matched entry wins. |


## Snippets

### Indexing unknown files

Currently, the way the dispatch algorithm works is if there is no matching handler,
e.g. say you don't have any tools to extract meaningful information from fonts,
those files won't even be indexed. This means that you won't be able to search for
them _at all_ using `xfs-search`. You can work around this by adding an additional
handler at the very end, e.g.:

```js
{
  "dispatch": [
    ["application/msword", "doc2text"],
    ["application/pdf",    "xfs-pdf"],
    ["application/x-text", "cat"],
    ["*", "echo"],
  ]
}
```
