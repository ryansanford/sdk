# Flywheel SDK

A golang SDK for interaction with a remote Flywheel instance.

## Building

```bash
git clone git@github.com:flywheel-io/core-sdk workspace/src/flywheel.io/sdk
ln -s workspace/src/flywheel.io/sdk sdk

./sdk/make.sh
```

## Testing

Three environment variables control the test suite:

* `SdkTestMode`: set to "unit" or "integration". Right now only integration is supported.
* `SdkTestHost`: In integration mode, set this to a host:port string. Defaults to "localhost:8443".
* `SdkTestKey`: In integration mode, set this to an API key. Defaults to "change-me".

If you're logged in with the flywheel CLI (`fw login`), running the integration suite is easy:

```bash
export SdkTestKey=$(jq -r .key ~/.config/flywheel/user.json)

./sdk/make.sh test
```

To run a single test:

```bash
./sdk/make.sh test -run TestSuite/TestGetConfig
```

## Route Implementation Status

Route                                                                                 | In SDK
--------------------------------------------------------------------------------------|--------
Get current user<br>`GET /users/self`                                                 | X
&nbsp;                                                                                |
Get all users<br>`GET /users`                                                         | X
Get user<br>`GET /users/<x>`                                                          | X
Add user<br>`POST /users/<x>`                                                         | X
Modify user<br>`PUT /users/<x>`                                                       | X
Delete user<br>`DELETE /users/<x>`                                                    | X
&nbsp;                                                                                |
Get all containers of type<br>`GET /<ctype>`                                          | X
Create container<br>`POST /<ctype>`                                                   | X
Get container<br>`GET /<ctype>/<x>`                                                   | X
Modify container<br>`PUT /<ctype>/<x>`                                                | X
Delete container<br>`DELETE /<ctype>/<x>`                                             | X
Upload file to container<br>`POST /<ctype>/<x>/files`                                 | X
Download file from container<br>`GET /<ctype>/<x>/files/<filename>`                   |
Get jobs that involve container<br>`GET /<ctype>/<x>/jobs`                            |
&nbsp;                                                                                |
Resolve path to route<br>`POST /resolve`                                              |
&nbsp;                                                                                |
Get all gears<br>`GET /gears`                                                         | X
Create gear<br>`POST /gears/<id>`                                                     | X
Get gear invocation<br>`GET /gears<x>/invocation`                                     | X
Suggest files for gear<br>`GET /gears<x>/suggest/<ctype>/<cid>`                       |
Delete gear<br>`DELETE /gears/<x>`                                                    | X
&nbsp;                                                                                |
Get a job<br>`GET /jobs/<x>`                                                          | X
Get a job's logs<br>`GET /jobs/<x>/logs`                                              | X
Append to a job's logs<br>`POST /jobs/<x>/logs`                                       | X
Enqueue a job<br>`POST /jobs/add`                                                     | X
Claim next pending job, and mark as running<br>`GET /jobs/next`                       | X
Modify job<br>`PUT /jobs/<x>`                                                         | X
&nbsp;                                                                                |
Get all batch jobs<br>`GET /batch`                                                    | X
Get batch job<br>`GET /batch/<x>`                                                     | X
Propose batch job<br>`POST /batch/<x>`                                                | X
Run batch job<br>`POST /batch/<x>/run`                                                | X
Cancel batch job<br>`POST /batch/<x>/cancel`                                          | X
&nbsp;                                                                                |
Create bulk download ticket<br>`POST /download`                                       |
Get bulk download from tricket<br>`GET /download`                                     |
&nbsp;                                                                                |
Various upload strategies?<br>`POST/upload/<label|uid|uid-match>`                     |
Engine upload<br>`POST /engine`                                                       |
&nbsp;                                                                                |
Declare a packfile upload to container<br>`POST /<ctype>/<x>/packfile-start`          |
Upload to packfile<br>`POST /<ctype>/<x>/packfile`                                    |
Complete packfile and listen for progress<br>`POST /<ctype>/<x>/packfile-end`         |
&nbsp;                                                                                |
Set project template<br>`POST /projects/<x>/template`                                 |
Delete project template<br>`DELETE /projects/<x>/template`                            |
Recalculate project template compliance<br>`POST /projects/<x>/recalc`                |
Recalculate all project template compliance<br>`POST /projects/recalc`                |
&nbsp;                                                                                |
Get all devices<br>`GET /devices`                                                     |
Get device<br>`GET /devices/<x>`                                                      |
Get device statuses<br>`GET /devices/status`                                          |
Get current device<br>`GET /devices/self`                                             |
&nbsp;                                                                                |
Get site configuration<br>`GET /config`                                               | X
Get site version<br>`GET /version`                                                    | X
&nbsp;                                                                                |
_Delayed_                                                                             |
Regenerate current user's API key<br>`POST /users/self/key`                           |
Get user groups<br>`GET /users/<x>/groups`                                            |
Get containers for user?<br>`GET /users/<x>/<ctype>`                                  |
Clean out expired packfile progress<br>`POST /clean-packfiles`                        |
Scan for and fix disconnected jobs<br>`POST /jobs/reap`                               |
Retry job<br>`POST /jobs/<x>/retry`                                                   |
List groups with projects the user can access<br>`GET /projects/groups`               |
Get schema<br>`GET /schemas/<schema>`                                                 |
Get job stats (redesign on the horizon)<br>`GET /jobs/stats`                          |
Get report<br>`GET /report/<site|project>`                                            |
Search?<br>`POST /search`                                                             |
Search files?<br>`GET /search/files`                                                  |
Search container?<br>`GET /search/<ctype>`                                            |
Get all gear rules<br>`GET /rules`                                                    |
Overwrite all gear rules<br>`POST /rules`                                             |
&nbsp;                                                                                |
_Won't be implemented_                                                                |
List known sites (depreciated)<br>`GET /sites`                                        |
Register a site (depreciated)<br>`POST /sites`                                        |
Get current user avatar (no point)<br>`GET /users/self/avatar`                        |
Get user avatar (no point)<br>`GET /users/<x>/avatar`                                 |
Get all jobs (depreciated) <br>`GET /jobs`                                            |
Get job configuration (no point)<br>`GET /jobs/<x>/config.json`                       |
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
