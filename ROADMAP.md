* [ ] Implement an install script.
    - [ ] Publish releases to GitHub when publishing a new tag.

* [ ] Add systemd support for easier service management on supported systems.

* [ ] Before releasing a non-alpha build explore what code can be exported
  to other projects, and what code should be flagged as internal.

* [ ] Support MongoDB as an alternative durable storage layer. Write
  concerns set to `majority` and read conerns set to `linearizable` should
  achieve a strongly consistent database that would suit this application.

* [ ] Add support for webhooks, to enable CI/CD and other workflows.

* [ ] Investigate storage improvements using BYTEA over CHAR(40) (this change
  would also enable support for SHA-256 Object IDs, more details on those IDs
  can be found here: <https://lwn.net/Articles/898522/>).
