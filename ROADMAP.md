* [ ] Implement an install script.
    - [ ] Publish releases to GitHub when publishing a new tag.

* [ ] Before releasing a non-alpha build explore what code can be exported
  to other projects, and what code should be flagged as internal.

* [ ] Support MongoDB as an alternative durable storage layer. Write
  concerns set to `majority` and read conerns set to `linearizable` should
  achieve a strongly consistent database that would suit this application.

