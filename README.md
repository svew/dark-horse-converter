# Dark Horse converter [![Build Status](https://travis-ci.org/Metalnem/dark-horse-converter.svg?branch=master)](https://travis-ci.org/Metalnem/dark-horse-converter) [![Go Report Card](https://goreportcard.com/badge/github.com/metalnem/dark-horse-converter)](https://goreportcard.com/report/github.com/metalnem/dark-horse-converter) [![license](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/metalnem/dark-horse-converter/master/LICENSE)
Converts digital Dark Horse comics into CBZ format. Comics are delivered to official iOS/Android clients as a tar archive that is not compatible with other comic book readers, so you have to convert them using this application before you can open them in the reader of your choice.

## Downloading comics

Easiest way to download Dark Horse comics is by using [Dark Horse downloader](https://chrome.google.com/webstore/detail/dark-horse-downloader/odciinkioeagogcogbpelccibomlenhl) Chrome extension.

If you don't want to install the extension, you can do it manually (following steps apply to Chrome web browser on Max OS X, but shouldn't be much different for any other OS/browser combination).

1. Open a web browser and go to [your bookshelf](https://digital.darkhorse.com/bookshelf).

2. Click on the comic book you want to convert.

3. View page source

    - In Chrome, click View -> Developer -> View source.
    - In Firefox, Ctrl+U

4. Search for the `bookreader-content` element in the page (Ctrl+F)

5. The line below should look something like this: `<a href="/books/9f1bb1e5dd524127bbdd0bba39d022e2">`. The long string of characters is the book uuid

6. Construct the tar archive link for the comic by replacing {uuid} part of the following string with the book_uuid that you found:

```
https://digital.darkhorse.com/api/v6/book/{uuid}
```

In this example, final link would look like this:

```
https://digital.darkhorse.com/api/v6/book/9f1bb1e5dd524127bbdd0bba39d022e2
```

8) Visit constructed link in your web browser to download the tar archive.

## Installation

```
$ go get github.com/metalnem/dark-horse-converter
```

## Binaries (x64)

[Windows](https://github.com/Metalnem/dark-horse-converter/releases/download/v1.0.0/dark-horse-converter-win64-1.0.0.zip)  
[Mac OS X](https://github.com/Metalnem/dark-horse-converter/releases/download/v1.0.0/dark-horse-converter-darwin64-1.0.0.zip)

## Usage

```
$ ./dark-horse-converter
Usage of ./dark-horse-converter:
  -i string
    	Path of a comic book file in tar format to convert (required)
  -o string
    	Output directory (if not specified, result will be placed in the same directory as the input)
```

## Example

```
$ ./dark-horse-converter -i 'The Witcher Sampler.tar' -o ~/Comics
```
