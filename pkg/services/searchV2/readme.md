Here is a proposal about a new package layout which:

* will make different types of indexes work separately, each consuming `entity_events` individually in its own control goroutine, so less contention
* will allow doing re-indexing with different periods for different types of indexes
* will allow doing isolated backups, will simplify extending in custom way

Entities:

```
type StandardSearchService struct {
  dashboardOrgIndexManager *orgIndexManager
}

// orgIndexManager is responsible controlling index lifecycle for all orgs: 
// re-indexing, making backups, polling entity_events (actually this is our 
// current `run` method).
type orgIndexManager struct {
  indexFactory func(ctx context.Context, orgID int64, writer *bluge.Writer) (Index, error)
  indexes map[int64]Index
}

func (*orgIndexManager) run(ctx) error {
  // Create initial indexes for known orgs.
  // Start applying events, do full re-indexing, do backups etc.
}

type Index interface {
  ReIndex(ctx) error
  ApplyUpdates(ctx, []EntityEvent) (newEventID, error)
  BackupTo(ctx, directory) error
}

type dashboardIndex struct {
  // Implements Index for one ORG.

  orgID int64
  writer *bluge.Writer
}
```

Backup directory is configurable per index type. Inside backup directory structure may look like this:

```
dashboard
   ├-- meta.json
   ├── org1
   │   ├── 0000000004a7.seg
   │   ├── 0000000004a8.seg
   │   ├── 0000000004a9.seg
   │   └── 000000000893.snp
   └── org2
```

`meta.json` contains:

```
{
  "eventId": 12
}
```

* `eventId` - the last event ID applied to the index in the backup. Since we make backups periodically we may need to apply some missing updates from `entity_event` table to catch up the state.

So possible flow may be like this:

1. event no.123
2. full reindex
3. backup - last event id 123
4. event no.124
5. event no.125
6. event no.126
7. node restart - restore backup and retrieve all events after 123

As we don't have org id separation in `entity_event` table we manage indexes for all organizations in one goroutine. By different types of indexes are a separate consumers of `entity_event` table - so different types of indexes do not depend on each other at all.

We can do backups after full-reindexing with the event id seen before re-indexing started.

Notes:

* At this moment we still need to re-index periodically since not all changes come to `entity_event` table.
* To reduce number of full re-indexes we can apply some checks – if we know that the database state matches the state in the index (i.e. no updates were missed) then we can skip re-index.
* We can also save/load backup to the remote storage, but need to preserve the structure.
* In HA scenario we can theoretically only have one node that does re-indexing, then share backup to all nodes - apply it and re-apply events from `entity_events`
* Re-indexing is still pretty resource greedy