#!/usr/bin/env python3
import io
import knita

kclient = knita.Client()

opts = [{"type": "host", "tags": {"name": "host-test"}},
        {"type": "docker", "docker_image": "ubuntu:latest", "docker_pull_strategy": "not-exists",
         "tags": {"name": "docker-test"}}]

for opts in opts:
    with kclient.runtime(**opts) as runtime:
        # Verify sysinfo is reported
        if runtime.sys_info() is None:
            raise Exception("sysinfo not set")
        if runtime.sys_info().os == "":
            raise Exception("sysinfo os not set")
        if runtime.sys_info().arch == "":
            raise Exception("sysinfo arch not set")
        if runtime.sys_info().total_cpu_cores <= 0:
            raise Exception("sysinfo cpu cores not set")
        if runtime.sys_info().total_memory <= 0:
            raise Exception("sysinfo memory not set")

        # Verify files can be imported
        expected_file_path = 'input/input.txt'
        with open(expected_file_path, "r") as file:
            expected_contents = file.read()
        runtime.import_(src=expected_file_path)
        runtime.exec(name="/bin/bash", args=["-c", f"contents=\"$(cat {expected_file_path})\"\n"
                                                   f"if [[ \"$contents\" != \"{expected_contents}\" ]]; then\n"
                                                   f"   exit 1\n"
                                                   f"fi\n"],
                     tags={"name": "import-test"})

        # Verify zero-byte files can be imported
        runtime.import_(src='input/zero-bytes.txt')
        runtime.exec(name="/bin/bash", args=["-c", f"stat input/zero-bytes.txt"],
                     tags={"name": "zero-byte-import-test"})

        # Verify the remote work directory is reported correctly
        runtime.exec(name="/bin/bash", args=["-c", f"contents=\"$(cat {runtime.work_directory(expected_file_path)})\"\n"
                                                   f"if [[ \"$contents\" != \"{expected_contents}\" ]]; then\n"
                                                   f"   exit 1\n"
                                                   f"fi\n"],
                     tags={"name": "work-directory-test"})

        # Verify files can be exported
        expected_contents = 'hello world\n'
        expected_file_path = 'output/host.txt'
        runtime.exec(name="/bin/bash",
                     args=["-c", f"mkdir output && echo -n '{expected_contents}' > {expected_file_path}"],
                     tags={"name": "export-test"})
        runtime.export(src=expected_file_path)
        with open(expected_file_path, "r") as file:
            contents = file.read()
        if contents != expected_contents:
            raise Exception("mismatched contents")

        # Verify stdout and stderr can be captured
        expected_output = 'hello world\n'
        with io.StringIO() as stdout, io.StringIO() as stderr:
            runtime.exec(name="/bin/bash", args=["-c", f" echo -n \"{expected_output}\" | tee /dev/stderr "],
                         tags={"name": "io-test"},
                         stdout=stdout,
                         stderr=stderr)
            if stdout.getvalue() != expected_output:
                raise Exception(f"mismatched stdout output: {stdout.getvalue()}")
            if stderr.getvalue() != expected_output:
                raise Exception(f"mismatched stderr output: {stderr.getvalue()}")

