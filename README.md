# Source

A horizontally scalable Git server.

## Introduction

This project aims to fully decouple the underlying durable storage behind
our Git repositories from the Git protocol itself, to reduce both the
technical and economical costs of horizontally scaling a Git server.

## Development status

Currently, Source supports all of the required behaviour needed to
collaborate successfully via Git.

Cloning, pushing, pulling, fetching, are all possible, with no noticeable
performance bottlenecks.

The web interface is limited, but is on par performance-wise with
Sourcehut, although lacks functionality beyond simply browsing
a repository.

## Future development

* Issue management via Git notes, with the ability to discuss and provide
  feedback within the repository itself.

* In repository support for patch-by-patch reviews, a highly successful
  approach to peer reviewing changes used often within Git mailing lists.

* Support partial checkouts without requiring changes to how Git stores
  objects internally.

## Current limitations

* Git LFS is not yet supported. Given that LFS is an extension to Git, and
  not core behaviour, implementing this is not currently a high-priority
  feature.

  That said, there's nothing stopping Source from supporting LFS. Many
  databases offer a range of options to support large file storage,
  including PostgreSQL.

