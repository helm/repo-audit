# Helm Repo Audit

This tool, as the name suggests, enables you to audit a Helm repository.

A tool like this could be run at regular intervals (k8s CronJob?) to audit a repo.

## TODO

* [ ] Check if digests for a given version have changed in a index file
* [ ] Email a report
* [ ] (Opt-in) check if newly released charts, since the last run, match their digest
* [ ] (opt-in) check if all charts match their digest (this will download all charts on all runs)
* [ ] (Opt-in) check and report which charts have provenance files
* [ ] (Opt-in) check new chart versions, since the last run, against their prov files
* [ ] (Opt-in) check all chart versions against their prov files
* [ ] Fire off a webhook on completion of an audit
* [ ] Audit metadata in the charts (need a plan for this)