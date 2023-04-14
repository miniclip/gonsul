# Gonsul - A Git to Consul tool, in Go

![Docker image](https://github.com/miniclip/gonsul/workflows/Docker%20image/badge.svg?branch=master)
![Checks and tests](https://github.com/miniclip/gonsul/workflows/Checks%20and%20tests/badge.svg?branch=master)

This tool serves as an entry point for the Hashicorp's Consul KV store. Not only because Consul
lacks of a built in
audit mechanism, but also because having configurations managed in GIT, using a gitflow or a normal
development-to-master flow is much friendly and familiar to any development team to manage
configurations.

Downloads in [releases page](https://github.com/miniclip/gonsul/releases).

## Important Notes

**The master branch is being refactored**. No (big) breaking changes will happen but codebase is
under significant changes.
Mainly adding tests, makefile, and code refactor/cleanup.

Please use the releases/downloads page.
Instructions for development will be added soon!

## How It Processes a Repository

Gonsul will (optionally) clone your repository into the filesystem. After, it will recursively parse
all the files in
the directory. Whenever Gonsul moves one level deep into a folder, the folder name is added as a
Consul KV path part
and as soon as it finds a file (either `.json`, `.txt` or `.ini` - or any other given in parameters)
it will take the
file name (without the extension) and append to the Consul KV path, making it a key and the file
content is added as the
value.

**Example:** Take this repository folder structure:

```plain
+-- prod
|   |
|   +-- website1
|   |   |-- config.json
|   |   |-- db-pass.txt
|   |
|   +-- app2
|       |-- config.json
|
+-- dev
    |
    +-- website1
    |   |-- config.json
    |
    +-- app2
        |-- config.json
```

All the files content would be inserted, respectively, into the following Consul KV paths:

- `prod/website1/config`
- `prod/website1/db-pass`
- `prod/app2/config`
- `dev/website1/config`
- `dev/app2/config`

**Note 1:** Gonsul makes the assumption that only `.txt`, `.json` and `.ini` files are to be treated
as keys and their
content as values for Consul. Any other files are disregarded. Also, `.json` files are treated and
validated as valid
JSON, and inserted into Consul pretty formatted.

**Note 2:** You can instruct Gonsul to expand the JSON files, so it will parse the structure,
appending the JSON path
to the previous folder structure, creating single Consul KV entries for each value. **Caveats:** Any
arrays found are
inserted into Consul a bracketed comma separated value string. More details on this on the flags
description.

## Some Features

Gonsul was built out of necessity and so it provides some features that makes our team more
comfortable regarding
security and consistency, so below are the most important features (TL;DR's in bold):

- The most obvious feature is it's built in GO, so we have **one small binary tool that just runs**
in any bamboo,
jenkins, *nix or OSX computer (does not support windows yet).
- Gonsul **uses Consul KV transactions**, so any operations (inserts, updates and deletes) will be
atomically done.
Due to Consul limitations, a maximum of 64 operations can be done on one transaction. If there are
more operations,
multiple batches are created and multiple calls made. Each transaction is atomic on Consul side, and
if a transaction
fails, Gonsul will terminate. On most configurations changes on a normal workday, all changes are
made using one atomic
transaction. Either all configs go, or none, reducing inconsistent states.
- Gonsul can **operate on a *less destructive* manner**. As we (@Miniclip) have multiple teams using
the KV store
without having the KV paths namespaced, we **strongly** rely on Consul ACL's. So to avoid a mass
delete in case of an
ACL mistake we would like to have Gonsul to operate on a mode where it does not perform any deletes.
In such mode, when
Gonsul detectes *delete* operation, it will terminate with a specific error number, printing a list
of KV paths that it
would delete, allowing a team member to decide whether to manually run the operation or investigate
any ACL problem
before disaster.
- Gonsul **handles secrets substitution** using [mustache.go](https://github.com/hoisie/mustache)
avoiding the need of
extra tooling for this small task. In our case, having all our configurations in GIT is a really
important way to keep a
familiar work flow, but having secrets out of GIT is crucial in any corporate environment. So just
use some placeholders
in your GIT files, and tell Gonsul where it can find a JSON key-value file of the
placeholder-secrets map. (More about
this below)
- Finally Gonsul **does not rely on any OS external tool** (such as GIT) as it does not shell out
any operations. All
tasks are natively handled by Gonsul, as long as filesystem and network permission/access exists.

## Available Flags

Below are all available command line flags for starting **Gonsul**. Flags may be specified on the
command line, using
environment variables (flag prefixed with "GONSUL_", converted to upper case, and "-" changed to
"_", e.g.
`GONSUL_REPO_URL`), or using a configuration file (`--config=`). The order of precedence for
specifying flags is:

1. Command line options
2. Environment variables
3. Configuration file
4. Default values

```bash
--config=
--strategy=
--repo-url=
--repo-ssh-key=
--repo-ssh-user=
--repo-branch=
--repo-remote-name=
--repo-base-path=
--repo-root=
--consul-url=
--consul-acl=
--consul-base-path=
--ca-file=
--key-file=
--cert-file=
--log-level=
--expand-json=
--secrets-file=
--allow-deletes=
--poll-interval=
--input-ext=
--keep-ext=
```

Below is the full description for each individual command line flag.

### `--config`

> `require:` **no**
> `example:` **`--config=gonsul.conf`**

This specifies a file with configuration settings for the options described below. File syntax:

- Empty lines and lines beginning with a "#" character are ignored.
- Flags and values can be separated with whitespace or the "=" character (`strategy ONCE` or
`strategy=DRYRUN`).
- Booleans can be empty (true) or set with "1/0", "t/f", "T/F", "true/false", "TRUE/FALSE", or
"True/False".

### `--strategy`

> `require:` **yes**
> `example:` **`--strategy=ONCE`**

This defines the mode that **Gonsul** should run at, having the following strategies:

- **`DRYRUN`** In this mode it will process the repository/folder and will **output only** the
different operations it
would run, such as: inserts, updates and deletes
- **`ONCE`** In this mode it will process the repository/folder and will proceed with all the
operations it finds
(updates, inserts and optionally deletes as well, see details below), after it finishes, Gonsul will
terminate
gracefully (exit 0). If something goes wrong, Gonsul will exit with one of the exit codes.
- **`POLL`** In this mode, Gonsul will start polling the repository every *X* minutes and check for
differences
**applying the same logic as the `once` strategy**. The process will work forever and it will only
stop on any error or
a *SIGINT* is received.
- **`HOOK`** This is the last strategy type, and in this mode, Gonsul will lauch an HTTP server that
listens on a
specified port waiting for `GET` request triggering and **applying the same logic as the `once`
strategy**. Again, this
process will live forever until it receives a *SIGINT* signal. Also, the HTTP server will only
process one request at a
time (using a lock between processes) to avoid concurrent writes on Consul, so whenever a request is
made and Gonsul is
already processing another, the new request will hold until the request before finishes.

**NOTES**: On both POLL and HOOK strategies, the application will gracefully terminate upon
receiving a `SIGINT` signal,
but it will obviously not under a `SIGKILL` in which you can end up with inconsistent data inside
Consul KV cluster (in
case you have more than 64 operations in progress, and Gonsul issued only the first transaction for
example)

### `--repo-url`

> `require:` **no**
> `example:` **`--repo-url=git@github.com:githubtraining/hellogitworld.git`**

This flag will tell Gonsul where it should clone/checkout the repository from. This should be a
fully qualified GIT url
(either ssh, http or https). Please provide the full URL, with any ports, credentials or whatever
you would normally use.

**Note:** If you do not provide this flag, Gonsul will still look at the filesystem `--repo-root`
folder and try to
process/parse the directory. This can be useful in case you want to delegate the repository cloning
to a CI platform,
such as Bamboo for example, removing the GIT responsibility from Gonsul and having it just process
the files and sync
them with Consul as normal.

### `--repo-ssh-key`

> `require:` **no**
> `example:` **`--repo-ssh-key=/home/example_user/.ssh/id_rsa`**

This is only required if the previous flag was given and it tells Gonsul where it can find the SSH
key file to use in
case we're connecting through SSH with a key.

### `--repo-ssh-user`

> `require:` **no**
> `example:` **`--repo-ssh-user=git`**

This is only required if the previous flag was given and it tell Gonsul what user we should use when
connecting to an
SSH GIT repository.

### `--repo-branch`

> `require:` **no**
> `default:` **master**
> `example:` **`--repo-branch=my_branch_name`**

This is the branch name that Gonsul should try to checkout.

### `--repo-remote-name`

> `require:` **no**
> `default:` **origin**
> `example:` **`--repo-remote-name=upstream_name`**

This is the name of the remote configured on the repository. We need it to properly trigger some
PULL and CHECKOUT
commands on the repository.

### `--repo-base-path`

> `require:` **no**
> `example:` **`--repo-base-path=configs/relative/path`**

This is the relative folder, inside our repository, that Gonsul should consider and parse as values
to be synced with
your Consul cluster. If no value is given, Gonsul assume that the root repository is to be
considered as Consul
parseable.

**Note:** Gonsul will only build the hierarchy path from this path onward. Given the example above,
Gonsul will not try
to look for a Consul KV path that starts with `configs/relative/path` but instead, it will from any
deeper path from
this folder down.

**Note2:** If this flag is not set, the below value of `--repo-root` will be used as the hierarchy
path within Consul,
which is normally not desired. Most of the times we should/must use this flag.

### `--repo-root`

> `require:` **yes**
> `example:` **`--repo-root=/home/user/gonsul/repo_dir`**

This is the absolute path where Gonsul should clone the repository in your OS filesystem. This is a
required field, and
when provided without providing the previous flag `--repo-url`, Gonsul will try to parse the given
absolute path without
using GIT. This could be useful if you want to use Gonsul as filesystem parser and Consul sync
without doing any GIT
operations.

**Note:** This value/path will be used as the the hierarchy path on Consul if no value is given on
above `--repo-base-path`
which in many cases is not intended. Most of the times, we should also use the flag above.

### `--consul-url`

> `require:` **yes**
> `example:` **`--consul-url=https://consul-cluster.example.com:8080`**

This is the Consul's REST API endpoint that Gonsul will call to make the required inserts, updates
and deletes on the
Consul's KV store. Please provide the full URL, with scheme and port if appropriate.

**Note:** Do not add the KV path, just the the scheme and the authority part of the URL.

### `--consul-acl`

> `require:` **no**
> `example:` **`--consul-acl=youracltokenhere`**

This is the Consul's access token that Gonsul will use when connecting to the cluster agent. This
ACL **must** have read
and write access to the same KV paths that are going to be replicated with the repository files.

### `--consul-base-path`

> `require:` **no**
> `example:` **`--consul-base-path=my/kv/base/path`**

This is the prefix for all generated keys that Gonsul will look at.

**Note:** Remember that this base path **must not** be mirrored in the repository, as it will be
automatically appended.
This is useful when the Consul cluster as all the KV paths segregated (namespaced) by teams or
projects.

### `--ca-file`

> `require:` **no**
> `example:` **`--ca-file=/path/my_ca.crt`**

This is the path for ca certfile for mTLS connection.

### `--cert-file`

> `require:` **no**
> `example:` **`--cert-file=/path/my_cert.crt`**

This is the path for cert certfile mTLS connection.

### `--key-file`

> `require:` **no**
> `example:` **`--key-file=/path/my_key.crt`**

This is the path for key certfile mTLS connection.

### `--log-level`

> `require:` **no**
> `default:` **ERROR**
> `example:` **`--log-level=DEBUG`**

This defines the logging level that **Gonsul** should run, having the following levels:

- `DEBUG` In this mode it will verbosely output the information of below levels plus all the
debugging information such
as what tasks are being processed.
- `INFO` In this mode it will output all the information below plus some useful information such
as processed keys.
- `ERROR` In this mode, Gonsul will only output only error messages.

### `--expand-json`

> `require:` **no**
> `default:` **false**
> `example:` **`--expand-json=true`**

This will tell Gonsul how to treat JSON files. If true, Gonsul will traverse the JSON files and
append the path to the
previous folder/file hierarchy and create single entries in Consul KV for each value.

**Note:** Because Consul is a simple KV Store, where **all values are strings**, there are some
caveats regarding the
JSON file expanding. Some important ones are:

- Any **arrays found** are inserted into Consul as a bracketed comma separated value string,
for example: `["val1","val2","otherval"]`
- All the **boolean values** are inserted into Consul as strings `true` and `false`. This might
break some applications
when reading configuration, as they will be just strings after all.
- All the **numeric values** will obviously be inserted into consul as strings. Again, take that
into consideration when
reading configurations from your app as any numeric values will be strings when coming out from
Consul.

### `--expand-yaml`

> `require:` **no**
> `default:` **false**
> `example:` **`--expand-yaml=true`**

This will tell Gonsul how to treat YAML files. If true, Gonsul will traverse the YAML files and
append the path to the
previous folder/file hierarchy and create single entries in Consul KV for each value.

**Note:** Because Consul is a simple KV Store, where **all values are strings**, there are some
caveats regarding the
YAML file expanding. Some important ones are:

- Any **arrays found** are inserted into Consul as a bracketed comma separated value string,
for example: `["val1","val2","otherval"]`
- All the **boolean values** are inserted into Consul as strings `true` and `false`. This might
break some applications
when reading configuration, as they will be just strings after all.
- All the **numeric values** will obviously be inserted into consul as strings. Again, take that
into consideration when
reading configurations from your app as any numeric values will be strings when coming out from
Consul.

### `--secrets-file`

> `require:` **no**
> `example:` **`--secret-files=secrets.json`**

This is the location for a secrets file. You can either pass a relative path and Gonsul will look
for it into the
`--repo-root` path or you can pass an absolute path. If this value is passed and the file is found
and checked as valid,
Gonsul will try to do an on-the-fly search and replace any placeholders for the corresponding
values. This is done using
[mustache](https://mustache.github.io/) template system.

This file must be a valid JSON map, where the keys are the placeholders and the values the actual
secrets.  An example
`secrets.json` file could look like this:

```json
{
   "FOO_DB_USER": "foo_username_foo1982",
   "FOO_DB_PASS": "foo_password_J5sXoEN",
   "BAR_TOKEN": "bar_token_J5sXoEN",
   "BAZ_KEY": "baz_key_J5sXoEN"
}
```

**Note 1:** All the replacement is done on-the-fly in memory, and apart from the original supplied
`secrets.json` file,
no secrets are written to disk.
**Note 2:** The placeholders **should** follow the *mustache* triple curly braces
`{{{FOO_DB_USER}}}`, that means
*"unescaped HTML charcaters"* - basically takes the value as is.

### `--allow-deletes`

> `require:` **no**
> `default:` **false**
> `example:` **`--allow-deletes=true`**

This option sets Gonsul behaviour on how it should proceed in case some Consul deletes are to be
made.

- **`true`** With this mode deletes are allowed and Gonsul will proceed with the delete operations,
removed files from the repo and/or files added directly on consul k/v will be deleted.
- **`skip`** In this mode, Gonsul will skip the deletion check and operations and will just
sync/push files to consul k/v path (can lead to inconsistencies at the consul k/v path since removed
files from the repo and/or files added directly on consul k/v will stay in the consul path because
we are skipping deletions).
- **`false`** *this is the default mode* and Gonsul will not proceed with the deletes, and depending
on the `--strategy` it is running at, will respond with some different behaviors, such as:

    1. **`ONCE`** When running in once mode, Gonsul will terminate with **error code 10** and output
    to console all the Consul KV paths that are supposed to be deleted.

    2. **`HOOK`** Gonsul will repond to the HTTP request with error 503 and will also return the
    following headers and values:
        - `X-Gonsul-Error:10`
        - `X-Gonsul-Delete-Paths:path1/to/be/deleted,path2/to/be/deleted`

    3. **`POLL`** Gonsul will log all the paths to be deleted as ERRORS and carry on, over and over.
    In this mode you should monitor Gonsul logs to detect any found errors, and react appropriately.
    The errors will follow the syntax:
        - `[ERROR] [28-01-2018 17:38:25 1234] error-10 path1/to/be/deleted`
        - `[ERROR] [28-01-2018 17:38:25 5678] error-10 path2/to/be/deleted`

### `--poll-interval`

> `require:` **no**
> `default:` **60**
> `example:` **`--poll-interval=300`**

This is the number of seconds you want Gonsul to wait between checks on the repository when it is
running in
`--strategy=POLL` mode.

### `--input-ext`

> `require:` **no**
> `default:` **json,txt,ini**
> `example:` **`--input-ext=json,txt,ini,yaml`**

This is the file extensions that Gonsul should consider as inputs to populate our Consul. Please
set each extension
without the dot, and separate each extension with a comma.

### `--keep-ext`

> `require:` **no**
> `default:` **false**
> `example:` **`--keep-ext=true`**

Gonsul default behavior is to remove/trim the filextension from the filename when pushing files to
consul k/v path, if this option is set to true gonsul will keep the file extension.

### `--timeout`

> `require:` **no**
> `default:` **5**
> `example:` **`--timeout=20`**

The number of seconds for the client to wait for a response from Consul

## Gonsul Exit Codes

Whenever an error occurs, and Gonsul exits with a code other than 0, we try to return a meaningful
code, such as:

- **10** - This is the most important error code. It means **Delete** operations were found
and Gonsul is running without
delete permission. This error comes with the info about the Consul KV paths that would be deleted.

- **20** - There was a problem on the initialization parameters /flags

- **30** - This means there was an error connecting to Consul cluster. This can ben either ACL
token, wrong endpoint, network, etc.

- **31** - There was a problem running a Consul transaction. It basically means on operation of
the transaction is corrupted for
some reason. Try a dryrun to analyze all the operations Gonsul is trying to run.

- **40** - This is a generic error when Gonsul fails to read an HTTP response.

- **50** - This error is thrown when Gonsul could not encode a json payload for a transaction.
Check **dryrun** for what operations
Gonsul is trying to run.

- **51** - This is when Gonsul could not decode a JSON payload. This can be either from a GET
response from Consul, or more common
when processing the filesystem and it found a corrupted JSON file - check your JSON files for
errors.

- **60** - This occurs when Gonsul cannot clone the repository. Either because credentials are
broken, or filesystem permissions.

- **70** - This error occurs when secret replacement fails.

- **80** - This is a generic HTTP error. Run Gonsul in debug mode to look for more information
regarding the error.

## Contributing

For notes on how to contribute check [CONTRIBUTING](CONTRIBUTING.md).
