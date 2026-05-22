import os
from typing import Optional, List
import grpc

from . import director_pb2
from . import executor_pb2
from . import director_pb2_grpc
from . import runtime


class ConfigException(Exception):
    """Raised when an invalid configuration value is provided to Knita."""
    pass


class Client:
    """Client connects back to the Knita CLI process to orchestrate builds."""
    __build_id: str
    __director_stub: director_pb2_grpc.DirectorStub

    def __init__(self):
        """Returns a Knita client that is configured to connect back to the Knita CLI process."""
        self.__build_id = os.environ.get('KNITA_BUILD_ID')
        if self.__build_id is None:
            raise ConfigException('expected KNITA_BUILD_ID to be set')

        socket = os.environ.get('KNITA_SOCKET')
        if socket is None:
            raise ConfigException('expected KNITA_SOCKET to be set')

        channel = grpc.insecure_channel(f'unix://{socket}')
        self.__director_stub = director_pb2_grpc.DirectorStub(channel)

    def runtime(self,
                type: runtime.RuntimeType,
                display_name: str = "",
                docker_image: str = None,
                docker_pull_strategy: runtime.DockerPullStrategy = None,
                docker_basic_auth: runtime.DockerBasicAuth = None,
                docker_aws_ecr_auth: runtime.DockerAWSECRAuth = None,
                labels: Optional[dict] = None,
                annotations: Optional[dict] = None,
                runs_on: Optional[dict] = None,
                runs_on_expressions: Optional[List[runtime.Requirement]] = None):
        """Opens a new remote runtime configured based on options.

        labels / annotations attach metadata to the runtime (carried through to its events).
        runs_on is a dict of matchLabels constraining which executor can host the runtime.
        runs_on_expressions is a list of label-selector Requirement entries (more expressive
        matching). Both runs_on and runs_on_expressions can be combined."""

        opts = executor_pb2.RuntimeOpts(
            display_name=display_name,
            meta=runtime._opts_meta(labels, annotations),
            label_selector=runtime._label_selector(runs_on, runs_on_expressions))
        if type == runtime.RuntimeType.host:
            opts.type = executor_pb2.RuntimeType.RUNTIME_HOST
        elif type == runtime.RuntimeType.docker:
            opts.type = executor_pb2.RuntimeType.RUNTIME_DOCKER
            opts.docker.image.image_uri = docker_image
            if docker_pull_strategy is not None:
                if docker_pull_strategy == runtime.DockerPullStrategy.always:
                    opts.docker.image.pull_strategy = executor_pb2.DockerPullOpts.PULL_STRATEGY_ALWAYS
                elif docker_pull_strategy == runtime.DockerPullStrategy.never:
                    opts.docker.image.pull_strategy = executor_pb2.DockerPullOpts.PULL_STRATEGY_NEVER
                elif docker_pull_strategy == runtime.DockerPullStrategy.not_exists:
                    opts.docker.image.pull_strategy = executor_pb2.DockerPullOpts.PULL_STRATEGY_NOT_EXISTS
                else:
                    raise ConfigException(f"Unknown Docker pull strategy: {docker_pull_strategy}")
            if docker_basic_auth is not None:
                opts.docker.image.auth = executor_pb2.BasicAuth(
                    username=docker_basic_auth.username, password=docker_basic_auth.password)
            if docker_aws_ecr_auth is not None:
                opts.docker.image.auth = executor_pb2.AWSECRAuth(
                    region=docker_aws_ecr_auth.region,
                    aws_access_key_id=docker_aws_ecr_auth.aws_access_key_id,
                    aws_secret_key=docker_aws_ecr_auth.aws_secret_key)
        else:
            raise ConfigException(f"Unknown runtime type: {type}")

        req = director_pb2.OpenRequest(build_id=self.__build_id, opts=opts)
        res = self.__director_stub.Open(req)
        return runtime.Runtime(res.runtime_id, res.work_directory, res.sys_info, self.__director_stub)
