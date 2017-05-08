#!/usr/bin/env python

import ctypes
import json
import sys
import os

# Load the shared object file. Further details are added at the end of the file
bridge = ctypes.cdll.LoadLibrary(os.path.join(os.path.dirname(__file__), '../c/flywheel.so'))

#
# Begin block to handle unicode in JSON
# http://stackoverflow.com/a/33571117
#

def _json_load_byteified(file_handle):
	return _byteify(
		json.load(file_handle, object_hook=_byteify),
		ignore_dicts=True
	)

def _json_loads_byteified(json_text):
	return _byteify(
		json.loads(json_text, object_hook=_byteify),
		ignore_dicts=True
	)

def _byteify(data, ignore_dicts = False):
	# if this is a unicode string, return its string representation
	if isinstance(data, unicode):
		return data.encode('utf-8')
	# if this is a list of values, return list of byteified values
	if isinstance(data, list):
		return [ _byteify(item, ignore_dicts=True) for item in data ]
	# if this is a dictionary, return dictionary of byteified keys and values
	# but only if we haven't already byteified it
	if isinstance(data, dict) and not ignore_dicts:
		return {
			_byteify(key, ignore_dicts=True): _byteify(value, ignore_dicts=True)
			for key, value in data.iteritems()
		}
	# if it's anything else, return it in its original form
	return data

#
# End block to handle unicode in JSON
#

def test_bridge(name):
	"""
	Test if the C bridge is functional.
	Should return "Hello <name>".
	"""

	pointer = bridge.TestBridge(name)
	payload = ctypes.cast(pointer, ctypes.c_char_p).value
	return payload

class FlywheelException(Exception):
	pass

class Flywheel:

	def __init__(self, key):
		splits = key.split(':')

		if len(splits) < 2:
			raise FlywheelException('Invalid API key.')

		self.key = key
		self.keyC = ctypes.create_string_buffer(key)

	@staticmethod
	def _handle_return(status, pointer):
		statusCode = status.value
		payload = ctypes.cast(pointer, ctypes.c_char_p).value

		if statusCode == 0 and payload is None:
			return None
		elif statusCode == 0:
			return _json_loads_byteified(payload)['data']
		else:
			result = 'Unknown error (status ' + str(statusCode) + ')'
			try:
				result = _json_loads_byteified(payload)['message']
			except:
				pass

			raise FlywheelException(result)

	#
	# AUTO GENERATED CODE FOLLOWS
	#

	{{range .Signatures}}
	def {{camel2snake .Name}}(self{{range .Params}}, {{.Name}}{{end}}):
		status = ctypes.c_int(-100)
		pointer = bridge.{{.Name}}(self.keyC, {{range .Params}}{{.Name}}, {{end}}ctypes.byref(status))
		return self._handle_return(status, pointer)
	{{end}}

# Every bridge function returns a char*.
# Manually informing ctypes of this prevents a segmentation fault on OSX.

# Manual functions
bridge.TestBridge.restype = ctypes.POINTER(ctypes.c_char)

# API client functions
{{range .Signatures}}bridge.{{.Name}}.restype = ctypes.POINTER(ctypes.c_char)
{{end}}
