#!/usr/bin/env python

from __future__ import print_function

import os
import random
import string

import flywheel

# Helper function to generate random data
def rand_string():
	return ''.join(random.choice(string.ascii_lowercase) for _ in range(10))


# Check that data can flow back & forth across the bridge
bridge_response = flywheel.test_bridge('world')
assert bridge_response == 'Hello world'

# Get API get from environment
api_key = os.environ['SdkTestKey']

# Create client
fw = flywheel.Flywheel(api_key)

# A test file
filepath = __file__
filename = os.path.basename(filepath)


#
## Users
#
user = fw.get_current_user()
assert len(user['_id']) > 0

users = fw.get_all_users()
assert len(users) >= 1

email = rand_string() + '@' + rand_string() + '.com'
userId = fw.add_user({
	'_id':       email,
	'email':     email,
	'firstname': rand_string(),
	'lastname':  rand_string(),
})

fw.modify_user(userId, {'firstname': 'John'})
user2 = fw.get_user(userId)
assert user2['email'] == email
assert user2['firstname'] == 'John'

fw.delete_user(userId)


#
## Groups
#
groupId = fw.add_group({
	'_id': rand_string(),
})

fw.add_group_tag(groupId, 'blue')
fw.modify_group(groupId, {'name': 'testdrive'})

groups = fw.get_all_groups()
assert len(groups) > 0

group = fw.get_group(groupId)
assert group['tags'][0] == 'blue'
assert group['name'] == 'testdrive'


#
## Projects
#
projectId = fw.add_project({
	'label': rand_string(),
	'group': groupId,
})

fw.add_project_tag(projectId, 'blue')
fw.modify_project(projectId, {'label': 'testdrive'})
fw.add_project_note(projectId, 'This is a note')

projects = fw.get_all_projects()
assert len(projects) > 0

fw.upload_file_to_project(projectId, filepath)
fw.download_file_from_project(projectId, filename, '/tmp/download.py')

project = fw.get_project(projectId)
assert project['tags'][0] == 'blue'
assert project['label'] == 'testdrive'
assert project['notes'][0]['text'] == 'This is a note'
assert project['files'][0]['name'] == filename
assert project['files'][0]['size'] == os.path.getsize('/tmp/download.py')


#
## Sessions
#
sessionId = fw.add_session({
	'label': rand_string(),
	'project': projectId,
})

fw.add_session_tag(sessionId, 'blue')
fw.modify_session(sessionId, {'label': 'testdrive'})
fw.add_session_note(sessionId, 'This is a note')

sessions = fw.get_project_sessions(projectId)
assert len(sessions) > 0

sessions = fw.get_all_sessions()
assert len(sessions) > 0

fw.upload_file_to_session(sessionId, filepath)
fw.download_file_from_session(sessionId, filename, '/tmp/download2.py')

session = fw.get_session(sessionId)
assert session['tags'][0] == 'blue'
assert session['label'] == 'testdrive'
assert session['notes'][0]['text'] == 'This is a note'
assert session['files'][0]['name'] == filename
assert session['files'][0]['size'] == os.path.getsize('/tmp/download2.py')


#
## Acquisitions
#
acqId = fw.add_acquisition({
	'label': rand_string(),
	'session': sessionId,
})

fw.add_acquisition_tag(acqId, 'blue')
fw.modify_acquisition(acqId, {'label': 'testdrive'})
fw.add_acquisition_note(acqId, 'This is a note')

acqs = fw.get_session_acquisitions(sessionId)
assert len(acqs) > 0

acqs = fw.get_all_acquisitions()
assert len(acqs) > 0

fw.upload_file_to_acquisition(acqId, filepath)
fw.download_file_from_acquisition(acqId, filename, '/tmp/download3.py')

acq = fw.get_acquisition(acqId)
assert acq['tags'][0] == 'blue'
assert acq['label'] == 'testdrive'
assert acq['notes'][0]['text'] == 'This is a note'
assert acq['files'][0]['name'] == filename
assert acq['files'][0]['size'] == os.path.getsize('/tmp/download3.py')


#
## Misc
#
config = fw.get_config()
assert config is not None

version = fw.get_version()
assert version['database'] >= 25


#
## Cleanup
#
fw.delete_acquisition(acqId)
fw.delete_session(sessionId)
fw.delete_project(projectId)
fw.delete_group(groupId)

print('')
print('Test drive complete.')