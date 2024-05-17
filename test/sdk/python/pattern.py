#!/usr/bin/env python3
import io
import knita

cc = knita.Client()

opts = [{"type": "host", "tags": {"name": "host-test"}},
        {"type": "docker", "docker_image": "ubuntu:latest", "docker_pull_strategy": "not-exists",
         "tags": {"name": "docker-test"}}]

for opts in opts:
    with cc.runtime(**opts) as runtime:
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
