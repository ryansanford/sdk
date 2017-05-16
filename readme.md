# Flywheel SDK

A golang SDK for interaction with a remote Flywheel instance.

[![Build status](https://circleci.com/gh/flywheel-io/core-sdk/tree/master.svg?style=shield)](https://circleci.com/gh/flywheel-io/core-sdk)

## Building

```bash
git clone https://github.com/flywheel-io/sdk workspace/src/flywheel.io/sdk
ln -s workspace/src/flywheel.io/sdk sdk

./sdk/make.sh
```

For other languages, check out the [bridge readme](bridge).

## Testing

Three environment variables control the test suite:

* `SdkTestMode`: set to `unit` or `integration`. Right now only integration is supported.
* `SdkTestKey`: In integration mode, set this to an API key. Defaults to `localhost:8443:change-me`.
* `SdkTestDebug`: Setting this to any value will cause each test to print an HTTP/1.1 representation of each request. Best used to debug a single failing test.

To run the integration test suite against a running API:

```bash
export SdkTestKey="localhost:8443:some-api-key"

./sdk/make.sh test
```

Or, to run a single test:

```bash
./sdk/make.sh test -run TestSuite/TestGetConfig
```

## Route Implementation Status

Route                                            | Golang  |  C++   | Python
-------------------------------------------------|---------|--------|--------
Get current user                                 | X       | X      | X
&nbsp;                                           |         |        |
Get all users                                    | X       | X      | X
Get user                                         | X       | X      | X
Add user                                         | X       | X      | X
Modify user                                      | X       | X      | X
Delete user                                      | X       | X      | X
&nbsp;                                           |         |        |
Get all containers of type                       | X       | X      | X
Create container                                 | X       | X      | X
Get container                                    | X       | X      | X
Modify container                                 | X       | X      | X
Delete container                                 | X       | X      | X
Upload file to container                         | X       | X      | X
Download file from container                     | X       | X      | X
Get jobs that involve container                  |         |        |
Support for collections                          |         |        |
Support for analyses                             |         |        |
&nbsp;                                           |         |        |
Resolve path to route                            |         |        |
&nbsp;                                           |         |        |
Get all gears                                    | X       | X      | X
Create gear                                      | X       | X      | X
Get gear invocation                              | X       |        |
Suggest files for gear                           |         |        |
Delete gear                                      | X       | X      | X
&nbsp;                                           |         |        |
Get a job                                        | X       | X      |
Get a job's logs                                 | X       | X      |
Append to a job's logs                           | X       |        |
Enqueue a job                                    | X       | X      |
Claim next pending job, and mark as running      | X       |        |
Modify job                                       | X       | X      | X
&nbsp;                                           |         |        |
Get all batch jobs                               | X       | X      |
Get batch job                                    | X       | X      |
Propose batch job                                | X       |        |
Start batch job                                  | X       | X      | X
Cancel batch job                                 | X       |        |
&nbsp;                                           |         |        |
Create bulk download ticket                      |         |        |
Get bulk download from tricket                   |         |        |
&nbsp;                                           |         |        |
Various upload strategies?                       |         |        |
Engine upload                                    |         |        |
&nbsp;                                           |         |        |
Declare a packfile upload to container           |         |        |
Upload to packfile                               |         |        |
Complete packfile and listen for progress        |         |        |
&nbsp;                                           |         |        |
Set project template                             |         |        |
Delete project template                          |         |        |
Recalculate project template compliance          |         |        |
Recalculate all project template compliance      |         |        |
&nbsp;                                           |         |        |
Get all devices                                  |         |        |
Get device                                       |         |        |
Get device statuses                              |         |        |
Get current device                               |         |        |
&nbsp;                                           |         |        |
Get site configuration                           | X       | X      | X
Get site version                                 | X       | X      | X
&nbsp;                                           |         |        |
_Delayed_                                        |         |        |
Regenerate current user's API key                |         |        |
Get user groups                                  |         |        |
Get containers for user?                         |         |        |
Clean out expired packfile progress              |         |        |
Scan for and fix disconnected jobs               |         |        |
Retry job                                        |         |        |
List groups with projects the user can access    |         |        |
Get schema                                       |         |        |
Get job stats (redesign on the horizon)          |         |        |
Get report                                       |         |        |
Search?                                          |         |        |
Search files?                                    |         |        |
Search container?                                |         |        |
Get all gear rules                               |         |        |
Overwrite all gear rules                         |         |        |
&nbsp;                                           |         |        |
_Won't be implemented_                           |         |        |
List known sites (depreciated)                   |         |        |
Register a site (depreciated)                    |         |        |
Get current user avatar (no point)               |         |        |
Get user avatar (no point)                       |         |        |
Get all jobs (depreciated)                       |         |        |
Get job configuration (no point)                 |         |        |

<!--

Left over for another day:


prefix('/<cont_name:{cname}>', [

	prefix('/<cid:{cid}>', [

		route('/<list_name:tags>',                 TagsListHandler,                     m=['POST']),
		route('/<list_name:tags>/<value:{tag}>',   TagsListHandler,                     m=['GET', 'PUT', 'DELETE']),

		route('/<list_name:files>',                FileListHandler,                     m=['POST']),
		route('/<list_name:files>/<name:{fname}>', FileListHandler,                     m=['GET', 'DELETE']),


		route('/<list_name:analyses>', AnalysesHandler, m=['POST']),
		# Could be in a prefix. Had weird syntax highlighting issues so leaving for another day
		route('/<list_name:analyses>/<_id:{cid}>',                       AnalysesHandler,                  m=['GET', 'DELETE']),
		route('/<list_name:analyses>/<_id:{cid}>/files',                 AnalysesHandler, h='download',    m=['GET']),
		route('/<list_name:analyses>/<_id:{cid}>/files/<name:{fname}>',  AnalysesHandler, h='download',    m=['GET']),
		route('/<list_name:analyses>/<_id:{cid}>/notes',                 AnalysesHandler, h='add_note',    m=['POST']),
		route('/<list_name:analyses>/<_id:{cid}>/notes/<note_id:{cid}>', AnalysesHandler, h='delete_note', m=['DELETE']),
		route('/<list_name:notes>',                                      NotesListHandler,                 m=['POST']),
		route('/<list_name:notes>/<_id:{nid}>',                          NotesListHandler, name='notes',   m=['GET', 'PUT', 'DELETE']),
	])
]),


prefix('/<cont_name:groups>', [
	route('/<cid:{gid}>/<list_name:roles>',                          ListHandler,     m=['POST']),
	route('/<cid:{gid}>/<list_name:roles>/<site:{sid}>/<_id:{uid}>', ListHandler,     m=['GET', 'PUT', 'DELETE']),

	route('/<cid:{gid}>/<list_name:tags>',                           TagsListHandler, m=['POST']),
	route('/<cid:{gid}>/<list_name:tags>/<value:{tag}>',             TagsListHandler, m=['GET', 'PUT', 'DELETE']),
]),




# Collections

route( '/collections',                 CollectionsHandler, h='get_all',                    m=['GET']),
route( '/collections',                 CollectionsHandler,                                 m=['POST']),
prefix('/collections', [
	route('/curators',                 CollectionsHandler, h='curators',                   m=['GET']),
	route('/<cid:{cid}>',              CollectionsHandler,                                 m=['GET', 'PUT', 'DELETE']),
	route('/<cid:{cid}>/sessions',     CollectionsHandler, h='get_sessions',               m=['GET']),
	route('/<cid:{cid}>/acquisitions', CollectionsHandler, h='get_acquisitions',           m=['GET']),
]),


# Collections / Projects

prefix('/<cont_name:collections|projects>', [
	prefix('/<cid:{cid}>', [
		route('/<list_name:permissions>',                          PermissionsListHandler, m=['POST']),
		route('/<list_name:permissions>/<site:{sid}>/<_id:{uid}>', PermissionsListHandler, m=['GET', 'PUT', 'DELETE']),
	]),
]),

# Misc (to be cleaned up later)

route('/<par_cont_name:groups>/<par_id:{gid}>/<cont_name:projects>', ContainerHandler, h='get_all', m=['GET']),
route('/<par_cont_name:{cname}>/<par_id:{cid}>/<cont_name:{cname}>', ContainerHandler, h='get_all', m=['GET']),

-->
