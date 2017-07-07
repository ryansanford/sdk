#!/usr/bin/env python

import ctypes
import json
import six
import sys
import os

# Load the shared object file. Further details are added at the end of the file.
bridge = ctypes.cdll.LoadLibrary(os.path.join(os.path.dirname(__file__), '../c/flywheelBridge.so'))

def test_bridge(s):
    """
    Test if the C bridge is functional.
    Should return "Hello <s>".
    """

    pointer = bridge.TestBridge(six.b(s))
    value = ctypes.cast(pointer, ctypes.c_char_p).value
    return value.decode('utf-8')

class FlywheelException(Exception):
    pass

class Flywheel:

    def __init__(self, key):
        if len(key.split(':')) < 2:
            raise FlywheelException('Invalid API key.')
        self.key = six.b(key)

    @staticmethod
    def _handle_return(status, pointer):
        status_code = status.value
        value = ctypes.cast(pointer, ctypes.c_char_p).value

        # In python 2, the casted pointer value will be of type str.
        # In python 3, it will instead be of type bytes.
        #
        # In python 3.6, json.loads gained the ability to process bytes objects.
        # Earlier versions did not have this capability.
        # So, to workaround, decode any non-str object.
        # This could later be changed to detect the python version -.-
        #
        # https://bugs.python.org/issue10976
        # https://bugs.python.org/msg275615
        if not isinstance(value, str):
            value = value.decode('utf-8')

        if status_code == 0 and value is None:
            return None
        elif status_code == 0:
            return json.loads(value)['data']
        else:
            try:
                msg = json.loads(value)['message']
            except:
                msg = 'Unknown error (status ' + str(status_code) + ').'
            raise FlywheelException(msg)

    #
    # AUTO GENERATED CODE FOLLOWS
    #
    {{range .Signatures}}
    def {{camel2snake .Name}}(self{{range .Params}}, {{.Name}}{{end}}):
        status = ctypes.c_int(-100)
        {{if ne .ParamDataName ""}}{{.ParamDataName}} = json.dumps({{.ParamDataName}})
        {{end -}}
        pointer = bridge.{{.Name}}(self.key, {{range .Params}}six.b(str({{.Name}})), {{end -}} ctypes.byref(status))
        return self._handle_return(status, pointer)
    {{end}}

# Every bridge function returns a char*.
# Declaring this explicitly prevents segmentation faults on OSX.

# Manual functions
bridge.TestBridge.restype = ctypes.POINTER(ctypes.c_char)

# API client functions
{{range .Signatures}}bridge.{{.Name}}.restype = ctypes.POINTER(ctypes.c_char)
{{end -}}
