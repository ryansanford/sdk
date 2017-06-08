%% Test methods in Bridge.m
%% Setup
disp('Setup')
% Before running this script, ensure the following paths were added
%   path to flywheel.so and flywheel.h files (same path)
%   path to Bridge.m
%   path to JSONlab

% Create string to be used in testdrive
testString = 'aeu8457bjclsj97v2h';
% A test file
currentFile = mfilename('fullpath');
filename = strcat(currentFile, '.m');
% Define error message
errMsg = 'Strings not equal';

% Create client
apiKey = getenv('SdkTestKey');
fw = Bridge(apiKey);

% Check that data can flow back & forth across the bridge
bridgeResponse = Bridge.testBridge('world');
assert(strcmp(bridgeResponse,'Hello world'), errMsg)
%% Users
disp('Testing Users')
user = fw.getCurrentUser();
assert(~isempty(user.x0x5F_id))

users = fw.getAllUsers();
assert(length(users) >= 1, 'No users returned')

% add a new user
email = strcat(testString, '@', testString, '.com');
userId = fw.addUser(struct('x0x5F_id',email,'email',email,'firstname',testString,'lastname',testString));

% modify the new user
fw.modifyUser(userId, struct('firstname', 'John'));
user2 = fw.getUser(userId);
assert(strcmp(user2.email, email), errMsg)
assert(strcmp(user2.firstname,'John'), errMsg)

fw.deleteUser(userId);

%% Groups
disp('Testing Groups')

groupId = fw.addGroup(struct('x0x5F_id',testString));

fw.addGroupTag(groupId, 'blue');
fw.modifyGroup(groupId, struct('name','testdrive'));

groups = fw.getAllGroups();
assert(~isempty(groups))

group = fw.getGroup(groupId);
assert(strcmp(group.tags,'blue'), errMsg)
assert(strcmp(group.name,'testdrive'), errMsg)

%% Projects
disp('Testing Projects')

projectId = fw.addProject(struct('label',testString,'group',groupId));

fw.addProjectTag(projectId, 'blue');
fw.modifyProject(projectId, struct('label','testdrive'));
fw.addProjectNote(projectId, 'This is a note');

projects = fw.getAllProjects();
assert(~isempty(projects), errMsg)

fw.uploadFileToProject(projectId, filename);
fw.downloadFileFromProject(projectId, filename, '/tmp/download.m');

project = fw.getProject(projectId);
assert(strcmp(project.tags,'blue'), errMsg)
assert(strcmp(project.label,'testdrive'), errMsg)
assert(strcmp(project.notes{1,1}.text, 'This is a note'), errMsg)
assert(strcmp(project.files{1,1}.name, filename), errMsg)
s = dir('/tmp/download.m');
assert(project.files{1,1}.size == s.bytes, errMsg)

%% Sessions
disp('Testing Sessions')

sessionId = fw.addSession(struct('label', testString, 'project', projectId));

fw.addSessionTag(sessionId, 'blue');
fw.modifySession(sessionId, struct('label', 'testdrive'));
fw.addSessionNote(sessionId, 'This is a note');

sessions = fw.getProjectSessions(projectId);
assert(~isempty(sessions), errMsg)

sessions = fw.getAllSessions();
assert(~isempty(sessions), errMsg)

fw.uploadFileToSession(sessionId, filename);
fw.downloadFileFromSession(sessionId, filename, '/tmp/download2.m');

session = fw.getSession(sessionId);
assert(strcmp(session.tags, 'blue'), errMsg)
assert(strcmp(session.label, 'testdrive'), errMsg)
assert(strcmp(session.notes{1,1}.text, 'This is a note'), errMsg)
assert(strcmp(session.files{1,1}.name, filename), errMsg)
s = dir('/tmp/download2.m');
assert(session.files{1,1}.size == s.bytes, errMsg)

%% Acquisitions
disp('Testing Acquisitions')

acqId = fw.addAcquisition(struct('label', testString,'session', sessionId));

fw.addAcquisitionTag(acqId, 'blue');
fw.modifyAcquisition(acqId, struct('label', 'testdrive'));
fw.addAcquisitionNote(acqId, 'This is a note');

acqs = fw.getSessionAcquisitions(sessionId);
assert(~isempty(acqs), errMsg)

%acqs = fw.getAllAcquisitions(); % TODO: couldn't load due to jsonlab not
% being able to load
%assert(~isempty(acqs), errMsg)

fw.uploadFileToAcquisition(acqId, filename);
fw.downloadFileFromAcquisition(acqId, filename, '/tmp/download3.m');

acq = fw.getAcquisition(acqId);
assert(strcmp(acq.tags,'blue'), errMsg)
assert(strcmp(acq.label,'testdrive'), errMsg)
assert(strcmp(acq.notes{1,1}.text, 'This is a note'), errMsg)
assert(strcmp(acq.files{1,1}.name, filename), errMsg)
s = dir('/tmp/download3.m');
assert(session.files{1,1}.size == s.bytes, errMsg)

%% Gears
disp('Testing Gears')

gearId = fw.addGear(struct('category','converter','exchange', struct('git_0x2D_commit','example','rootfs_0x2D_hash','sha384:example','rootfs_0x2D_url','https://example.example'),'gear', struct('name','test-drive-gear','label','Test Drive Gear','version','3','author','None','description','An empty example gear','license','Other','source','http://example.example','url','http://example.example','inputs', struct('x', struct('base','file')))));

gear = fw.getGear(gearId);
assert(strcmp(gear.gear.name, 'test-drive-gear'), errMsg)

gears = fw.getAllGears();
assert(~isempty(gears), errMsg)

job2Add = struct('gear_id',gearId,'state','pending','inputs',struct('x',struct('type','acquisition','id',acqId,'name','testDrive.m')));
jobId = fw.addJob(job2Add);

job = fw.getJob(jobId);
assert(strcmp(job.gear_id,gearId), errMsg)

logs = fw.getJobLogs(jobId);
% Likely will not have anything in them yet

%% Misc
disp('Testing Misc')

config = fw.getConfig();
assert(~isempty(config), errMsg)

version = fw.getVersion();
assert(version.database >= 25, errMsg)

%% Cleanup
disp('Cleanup')

fw.deleteAcquisition(acqId);
fw.deleteSession(sessionId);
fw.deleteProject(projectId);
fw.deleteGroup(groupId);
fw.deleteGear(gearId);

disp('')
disp('Test drive complete.')
