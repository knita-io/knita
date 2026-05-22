import os
from enum import Enum
from typing import Optional, List
from google.protobuf.any_pb2 import Any
from . import director_pb2
from . import director_pb2_grpc
from . import executor_pb2
from . import builtin_pb2


class Operator(str, Enum):
    """Operator values for a label-selector Requirement (matches the gRPC proto)."""
    in_ = "in"
    not_in = "not-in"
    exists = "exists"
    does_not_exist = "not-exists"


class Requirement:
    """A single matchExpression entry on a LabelSelector."""

    def __init__(self, key: str, operator: Operator, values: Optional[List[str]] = None):
        self.key = key
        self.operator = operator
        self.values = values or []


def _opts_meta(labels: Optional[dict] = None,
               annotations: Optional[dict] = None) -> Optional[executor_pb2.OptsMeta]:
    """Build an OptsMeta from labels/annotations, or return None if both are empty."""
    if not labels and not annotations:
        return None
    return executor_pb2.OptsMeta(labels=labels or {}, annotations=annotations or {})


def _label_selector(match_labels: Optional[dict],
                    match_expressions: Optional[List[Requirement]]) -> Optional[executor_pb2.LabelSelector]:
    """Build a LabelSelector from a matchLabels dict and a list of Requirement, or None."""
    if not match_labels and not match_expressions:
        return None
    op_map = {
        Operator.in_: executor_pb2.LabelSelectorRequirement.IN,
        Operator.not_in: executor_pb2.LabelSelectorRequirement.NOT_IN,
        Operator.exists: executor_pb2.LabelSelectorRequirement.EXISTS,
        Operator.does_not_exist: executor_pb2.LabelSelectorRequirement.DOES_NOT_EXIST,
    }
    exprs = []
    for r in match_expressions or []:
        exprs.append(executor_pb2.LabelSelectorRequirement(
            key=r.key, operator=op_map[r.operator], values=r.values))
    return executor_pb2.LabelSelector(matchLabels=match_labels or {}, matchExpressions=exprs)


class RuntimeType(str, Enum):
    """
    An enum specifying the type of runtime to create.
    Attributes:
        host: Host is runtime that executes directly on the host executor without any containerization or virtualization.
        docker: Docker is a runtime that executes inside a Docker container.
    """
    host = 'host'
    docker = 'docker'


class DockerPullStrategy(str, Enum):
    """
    An enum specifying when to pull a Docker image during runtime creation.
    Attributes:
        always: Always ensures that a Docker pull is performed prior to starting the runtime.
        never: Never will skip the Docker pull entirely. This is useful if executors are loaded with relevant Docker
               images out of band of any individual build.
        not_exists: NotExists will run Docker pull only if no matching image is found on the executor.
    """
    always = 'always'
    never = 'never'
    not_exists = 'not-exists'


class DockerBasicAuth:
    username: str
    password: str

    def __init__(self, username: str, password: str):
        self.username = username
        self.password = password


class DockerAWSECRAuth:
    region: str
    aws_access_key_id: str
    aws_secret_key: str

    def __init__(self, region: str, aws_access_key_id: str, aws_secret_key: str):
        self.region = region
        self.aws_access_key_id = aws_access_key_id
        self.aws_secret_key = aws_secret_key


class ExecException(Exception):
    """Raised when exec finishes with a non-zero exit code."""
    exit_code: int

    def __init__(self, exit_code: int):
        self.exit_code = exit_code


class Runtime:
    """Runtime represents a local handle to a remote runtime hosted by an executor.
    """
    __runtime_id: str
    __remote_work_directory: str
    __remote_sys_info: executor_pb2.SystemInfo
    __director_stub: director_pb2_grpc.DirectorStub

    def __init__(self, runtime_id: str, remote_work_directory: str, remote_sys_info: executor_pb2.SystemInfo,
                 director_stub: director_pb2_grpc.DirectorStub):
        self.__runtime_id = runtime_id
        self.__remote_work_directory = remote_work_directory
        self.__remote_sys_info = remote_sys_info
        self.__director_stub = director_stub

    def __enter__(self):
        return self

    def __exit__(self, *args):
        self.close()

    def id(self) -> str:
        """Returns the unique ID of the runtime."""
        return self.__runtime_id

    def sys_info(self) -> executor_pb2.SystemInfo:
        """Returns information about the runtime execution environment."""
        return self.__remote_sys_info

    def work_directory(self, rel_path: str = None) -> str:
        """Returns the fully qualified remote work directory of the runtime.
        Specify a relative path to have it joined to the work directory.
        This is helpful when exec'ing commands that reference file paths within the runtime."""
        if rel_path is None:
            return self.__remote_work_directory
        else:
            return os.path.join(self.__remote_work_directory, rel_path)

    def import_(self, src: str, dest: str = None, excludes: [str] = None,
                display_name: str = "", labels: Optional[dict] = None,
                annotations: Optional[dict] = None):
        """Import files from the local work directory into the runtime's remote work directory. src must be a
        relative path, and may be a glob (doublestar syntax supported). By default, all files identified by src will
        be copied to their original location on the remote. Use ImportOpts.dest to override this."""
        req = director_pb2.ImportRequest(
            runtime_id=self.__runtime_id,
            opts=director_pb2.ImportOpts(src_path=src, dest_path=dest, excludes=excludes,
                                         display_name=display_name,
                                         meta=_opts_meta(labels, annotations)))
        self.__director_stub.Import(req)

    def export(self, src: str, dest: str = None, excludes: [str] = None,
               display_name: str = "", labels: Optional[dict] = None,
               annotations: Optional[dict] = None):
        """Export files from the runtime's remote work directory into the local work directory. src must be a
        relative path, and may be a glob (doublestar syntax supported). By default, all files identified by src will
        be copied to their original location locally. Use ExportOpts.dest to override this."""
        req = director_pb2.ExportRequest(
            runtime_id=self.__runtime_id,
            opts=director_pb2.ExportOpts(src_path=src, dest_path=dest, excludes=excludes,
                                         display_name=display_name,
                                         meta=_opts_meta(labels, annotations)))
        self.__director_stub.Export(req)

    def exec(self, name: str, args: [str] = None, env: [str] = None, display_name: str = "", stdout=None,
             stderr=None, labels: Optional[dict] = None, annotations: Optional[dict] = None):
        """Exec executes a command inside the remote runtime.
        Raises ExecException if the command finishes with a non-zero code."""
        req = director_pb2.ExecRequest(
            runtime_id=self.__runtime_id,
            opts=executor_pb2.ExecOpts(name=name, args=args, env=env,
                                       display_name=display_name,
                                       meta=_opts_meta(labels, annotations)))
        for event in self.__director_stub.Exec(req):
            any_msg: Any = event.payload
            type_name = any_msg.TypeName()
            if type_name == builtin_pb2.StdoutEvent.DESCRIPTOR.full_name:
                ev = builtin_pb2.StdoutEvent()
                any_msg.Unpack(ev)
                if stdout:
                    stdout.write(ev.data.decode())
                continue
            if type_name == builtin_pb2.StderrEvent.DESCRIPTOR.full_name:
                ev = builtin_pb2.StderrEvent()
                any_msg.Unpack(ev)
                if stderr:
                    stderr.write(ev.data.decode())
                continue
            if type_name == builtin_pb2.ExecEndEvent.DESCRIPTOR.full_name:
                ev = builtin_pb2.ExecEndEvent()
                any_msg.Unpack(ev)
                if ev.HasField("result"):
                    code = ev.result.exit_code
                    if code != 0:
                        raise ExecException(code)
                    return

    def close(self):
        """Close the runtime. After a call to close the runtime can no longer be used."""
        req = executor_pb2.CloseRequest(runtime_id=self.__runtime_id)
        self.__director_stub.Close(req)
