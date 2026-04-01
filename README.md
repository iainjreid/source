# Source

Source is a horizontally scalable, _serverless_ Git server.

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

## Current limitations

* Git LFS is not yet supported. Given that LFS is an extension to Git, and
  not core behaviour, implementing this is not currently a high-priority
  feature.

  That said, there's nothing stopping Source from supporting LFS. Many
  databases offer a range of options to support large file storage,
  including PostgreSQL.

