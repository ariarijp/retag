# Retag

Exporting and Importing tool for Redash Query name and tags.

## Installation

```
go install github.com/ariarijp/retag@latest
```

### Setup

```
export REDASH_API_KEY="" # You can only use Admin user's API Key when you want to use import command.
export REDASH_URL="http://localhost:5000" # Without trailing slash
```

### Export queries

```
retag export --url $REDASH_URL --api-key $REDASH_API_KEY > example.yml
cat example.yml
queries:
  - id: 1
    name: Query 1
    tags: []
  - id: 2
    name: Query 2
    tags: [foo, bar]
  - id: 3
    name: Query 3
    tags: [baz]
```

### Edit YAML file

```yaml
queries:
  - id: 1
    name: Query 1
    tags: [ ]
  - id: 2
    name: Query 2
    tags: [ foo, bar ]
  - id: 3
    name: Query 3
    tags: [ baz ]
```

### Show diff

```
retag diff --url $REDASH_URL --api-key $REDASH_API_KEY --yaml example.yml
 queries:
   - id: 1
     name: Query 1
-    tags: []
+    tags: [foo, bar, baz]
   - id: 2
-    name: Query 2
+    name: Great Query
     tags: [foo, bar]
   - id: 3
-    name: Query 3
-    tags: [baz]
+    name: Awesome Query
+    tags: [bar, baz]

```

### Import queries

#### Dry run

```
retag import --url $REDASH_URL --api-key $REDASH_API_KEY --yaml example.yml # Dry run
DRYRUN: Query #1 will be updated.
DRYRUN: Query #2 will be updated.
DRYRUN: Query #3 will be updated.
```

#### Run

```
retag import --url $REDASH_URL --api-key $REDASH_API_KEY --yaml example.yml --dry-run=false
RESULT: Query #1 updated.
RESULT: Query #2 updated.
RESULT: Query #3 updated.
```

### Show diff again(has no changes)

```
retag diff --url $REDASH_URL --api-key $REDASH_API_KEY --yaml example.yml                
 queries:
   - id: 1
     name: Query 1
     tags: [foo, bar, baz]
   - id: 2
     name: Great Query
     tags: [foo, bar]
   - id: 3
     name: Awesome Query
     tags: [bar, baz]

```