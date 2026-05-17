# Postgres

## Bulk insert strategy

* Use COPY command across multiple connections to parallelise the
  consumption of data within Postgres. 

* Avoid parallel writes for smaller write batches, to prevent push events
  from single machines overconsuming database resources.
