import os
from enum import Enum

from . import director_pb2
from . import director_pb2_grpc
from . import executor_pb2


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

    def __init__(self, runtime_id: str, remote_work_directory: str, remote_sys_info: executor_pb2.SystemInfo, director_stub: director_pb2_grpc.DirectorStub):
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

    def import_(self, src: str, dest: str = ""):
        """Import files from the local work directory into the runtime's remote work directory.
        src and dest must be relative paths. src may be a glob (doublestar syntax supported).
        If dest is empty, all files identified by src will be copied to their original location in dest."""
        req = director_pb2.ImportRequest(runtime_id=self.__runtime_id, src_path=src, dest_path=dest)
        self.__director_stub.Import(req)

    def export(self, src: str, dest: str = ""):
        """Export files from the runtime's remote work directory into the local work directory.
        src and dest must be relative paths. src may be a glob (doublestar syntax supported).
        If dest is empty, all files identified by src will be copied to their original location in dest."""
        req = director_pb2.ExportRequest(runtime_id=self.__runtime_id, src_path=src, dest_path=dest)
        self.__director_stub.Export(req)

    def exec(self, name: str, args: [str] = None, env: [str] = None, tags: dict[str, str] = None, stdout=None,
             stderr=None):
        """Exec executes a command inside the remote runtime.
        Check the returned ExecResponse to see the command's exit code (a non-zero code is not an exception)."""
        req = director_pb2.ExecRequest(runtime_id=self.__runtime_id,
                                       opts=executor_pb2.ExecOpts(name=name, args=args, env=env, tags=tags))
        for event in self.__director_stub.Exec(req):
            field = event.WhichOneof('payload')
            payload = getattr(event, field)
            if field == 'exec_end':
                if payload.exit_code is not 0:
                    raise ExecException(payload.exit_code)
                return
            elif field == 'stdout':
                if stdout is not None:
                    stdout.write(payload.data.decode())
            elif field == 'stderr':
                if stderr is not None:
                    stderr.write(payload.data.decode())

    def close(self):
        """Close the runtime. After a call to close the runtime can no longer be used."""
        req = executor_pb2.CloseRequest(runtime_id=self.__runtime_id)
        self.__director_stub.Close(req)
