#!/usr/bin/env python

import ctypes
import json
import sys
import os

# Load the shared object file. Further details are added at the end of the file
bridge = ctypes.cdll.LoadLibrary(os.path.join(os.path.dirname(__file__), '../c/flywheel.so'))

def test_bridge(name):
    """
    Test if the C bridge is functional.
    Should return "Hello <name>".
    """

    pointer = bridge.TestBridge(bytes(name, 'utf-8'))
    payload = ctypes.cast(pointer, ctypes.c_char_p).value
    return payload.decode('utf-8')

class FlywheelException(Exception):
    pass

class Flywheel:

    def __init__(self, key):
        splits = key.split(':')

        if len(splits) < 2:
            raise FlywheelException('Invalid API key.')

        self.key = key
        self.keyC = ctypes.create_string_buffer(bytes(key, 'utf-8'))

    @staticmethod
    def _handle_return(status, pointer):
        statusCode = status.value
        payload = ctypes.cast(pointer, ctypes.c_char_p).value

        if statusCode == 0 and payload is None:
            return None
        elif statusCode == 0:
            return json.loads(payload)['data']
        else:
            result = 'Unknown error (status ' + str(statusCode) + ')'
            try:
                result = json.loads(payload)['message']
            except:
                pass

            raise FlywheelException(result)

    #
    # AUTO GENERATED CODE FOLLOWS
    #

    {{range .Signatures}}
    def {{camel2snake .Name}}(self{{range .Params}}, {{.Name}}{{end}}):
        status = ctypes.c_int(-100)
        {{if ne .ParamDataName ""}}marshalled_{{.ParamDataName}} = json.dumps({{.ParamDataName}})
        {{end}}
        pointer = bridge.{{.Name}}(self.keyC, {{range .Params}}bytes({{if eq .Type "data"}}marshalled_{{end}}{{.Name}}, 'utf-8'), {{end}}ctypes.byref(status))
        return self._handle_return(status, pointer)
    {{end}}

# Every bridge function returns a char*.
# Manually informing ctypes of this prevents a segmentation fault on OSX.

# Manual functions
bridge.TestBridge.restype = ctypes.POINTER(ctypes.c_char)

# API client functions
{{range .Signatures}}bridge.{{.Name}}.restype = ctypes.POINTER(ctypes.c_char)
{{end}}
