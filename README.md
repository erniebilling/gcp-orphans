# gcp-orphans
Tool to discover and delete orphaned GCP resources

A work in progress.

```
git clone https://github.com/erniebilling/gcp-orphans.git
cd gcp-orphans
make
```

Currently discovers orphaned firewall rules
```
./gcp-orphans firewallrules -g <path to service account JSON key file>
```

