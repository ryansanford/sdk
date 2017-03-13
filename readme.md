# Flywheel SDK

A golang SDK for interaction with a remote Flywheel instance.

## Route Implementation Status

List from routing table in [api.py](https://github.com/scitran/core/blob/master/api/api.py):

Route                                                                                 | In SDK
--------------------------------------------------------------------------------------|--------
Get current user<br>`GET /users/self`                                                 | X
Regenerate current user's API key<br>`POST /users/self/key`                           |
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
Get jobs that involve container<br>`GET /<ctype>/<x>/jobs`                            |
&nbsp;                                                                                |
Resolve path to route<br>`POST /resolve`                                              |
&nbsp;                                                                                |
Create bulk download ticket<br>`POST /download`                                       |
Get bulk download from tricket<br>`GET /download`                                     |
&nbsp;                                                                                |
Various upload strategies?<br>`POST/upload/<label|uid|uid-match>`                     |
Engine upload<br>`POST /engine`                                                       |
&nbsp;                                                                                |
Get all jobs<br>`GET /jobs`                                                           |
Claim next pending job, and mark as running<br>`GET /jobs/next`                       |
Enqueue a job<br>`POST /jobs/add`                                                     |
Get job<br>`GET /jobs/<x>`                                                            |
Modify job<br>`PUT /jobs/<x>`                                                         |
Get job configuration<br>`GET /jobs/<x>/config.json`                                  |
Retry job<br>`POST /jobs/<x>/retry`                                                   |
Get job stats<br>`GET /jobs/stats`                                                    |
Get all gears<br>`GET /gears`                                                         |
Get gear invocation<br>`GET /gears<x>/invocation`                                     |
Suggest files for gear<br>`GET /gears<x>/suggest/<ctype>/<cid>`                       |
&nbsp;                                                                                |
Get all gear rules<br>`GET /rules`                                                    |
Overwrite all gear rules<br>`POST /rules`                                             |
&nbsp;                                                                                |
Get all batch jobs<br>`GET /batch`                                                    |
Get batch job<br>`GET /batch/<x>`                                                     |
Run batch job<br>`POST /batch/<x>/run`                                                |
Cancel batch job<br>`POST /batch/<x>/cancel`                                          |
Get jobs from batch<br>`GET /batch/<x>/jobs`                                          |
&nbsp;                                                                                |
Declare a packfile upload to container<br>`POST /<ctype>/<x>/packfile-start`          |
Upload to packfile<br>`POST /<ctype>/<x>/packfile`                                    |
Complete packfile and listen for progress<br>`POST /<ctype>/<x>/packfile-end`         |
&nbsp;                                                                                |
Get report<br>`GET /report/<site|project>`                                            |
Search?<br>`POST /search`                                                             |
Search files?<br>`GET /search/files`                                                  |
Search container?<br>`GET /search/<ctype>`                                            |
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
Get config<br>`GET /config`                                                           | X
Get version<br>`GET /version`                                                         | X
&nbsp;                                                                                |
_Delayed_                                                                             |
Get user groups<br>`GET /users/<x>/groups`                                            |
Get containers for user?<br>`GET /users/<x>/<ctype>`                                  |
Clean out expired packfile progress<br>`POST /clean-packfiles`                        |
Scan for and fix disconnected jobs<br>`POST /jobs/reap`                               |
List groups with projects the user can access<br>`GET /projects/groups`               |
Get schema<br>`GET /schemas/<schema>`                                                 |
&nbsp;                                                                                |
_Won't be implemented_                                                                |
List known sites (depreciated)<br>`GET /sites`                                        |
Register a site (depreciated)<br>`POST /sites`                                        |
Get current user avatar (no point)<br>`GET /users/self/avatar`                        |
Get user avatar (no point)<br>`GET /users/<x>/avatar`                                 |

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
