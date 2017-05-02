
	def {{.Name}}(self{{range .Params}}, {{.Name}}{{end}}):
		status = ctypes.c_int(-100)
		pointer = bridge.{{.Name}}(self.keyC, {{range .Params}}{{.Name}}, {{end}}ctypes.byref(status))
		payload = ctypes.c_char_p(pointer).value

		return self._handle_return(status, payload)
