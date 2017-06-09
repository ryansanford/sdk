% Flywheel
classdef Flywheel
    properties
        key
    end
    methods
        function obj = Flywheel(apiKey)
            % Check if key is valid
            %   key must be in format <domain>:<API token>
            C = strsplit(apiKey, ':');
            if length(C) < 2
                ME = MException('FlywheelException:Invalid', 'Invalid API Key');
                throw(ME)
            end
            obj.key = apiKey;

            % Load flywheel shared library
            if not(libisloaded('flywheelBridge'))
                % loading in flywheelBridge.so file
                loadlibrary('flywheelBridge','flywheelBridgeSimple.h')
            end

            % Suppress Max Length Warning
            warningid = 'MATLAB:namelengthmaxexceeded';
            warning('off',warningid);
        end
        %
        % AUTO GENERATED CODE FOLLOWS
        %

        {{range .Signatures}}% {{camel2lowercamel .Name}}
        function result = {{camel2lowercamel .Name}}(obj{{range .Params}}, {{.Name}}{{end}})
            statusPtr = libpointer('int32Ptr',-100);
            {{if ne .ParamDataName ""}}{{.ParamDataName}} = savejson('',{{.ParamDataName}});
            {{end -}}
            pointer = calllib('flywheelBridge','{{.Name}}',obj.key,{{range .Params}}{{.Name}},{{end -}} statusPtr);
            result = Flywheel.handleJson(statusPtr,pointer);
        end
        {{end}}
        % AUTO GENERATED CODE ENDS
    end
    methods (Static)
        % Handle JSON using JSONlab
        function structFromJson = handleJson(statusPtr,ptrValue)
            % Get status value
            statusValue = statusPtr.Value;
            % If status indicates success, load JSON
            if statusValue == 0
                % Interpret JSON string blob as a struct object
                loadedJson = loadjson(ptrValue);
                % loadedJson contains status, message and data, only return
                %   the data information
                structFromJson = loadedJson.data;
            % Otherwise, nonzero statusCode indicates an error
            else
                % Try to load message from the JSON
                try
                    loadedJson = loadjson(ptrValue);
                    msg = loadedJson.message;
                    ME = MException('FlywheelException:handleJson', msg);
                % If unable to load message, throw an 'unknown' error
                catch ME
                    msg = sprintf('Unknown error (status %d).',statusValue);
                    causeException = MException('FlywheelException:handleJson', msg);
                    ME = addCause(ME,causeException);
                    rethrow(ME)
                end
                throw(ME)
            end
        end
        % TestBridge
        function ptrValue = testBridge(s)
            % Call bridge
            ptrValue = calllib('flywheelBridge','TestBridge',s);
        end
    end
end
