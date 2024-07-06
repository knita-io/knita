# Getting Started with the Knita SDK for Python

In this guide we'll create a simple "hello world" build using the Knita SDK for Python. In order to run this, you will
need
Docker installed.

1. Download the latest Knita CLI from the [release page](https://github.com/knita-io/knita/releases) and make sure it's
   in your path

2. Create a new `pattern.py` Python file and copy in the example below:

   ```python
   #!/usr/bin/env python3
   import knita
   
   client = knita.Client()
   
   with client.runtime(type='docker', docker_image='alpine:latest', tags={'name': 'example'}) as runtime:
       runtime.exec(name='/bin/sh', args=['-c', 'echo "hello world"'], tags={'name': 'hello-world'})
   ```

   ```bash
   chmod +x pattern.py
   ```

3. Install the Knita Python package into a virtual env:

   ```bash
   python3 -m venv venv
   source venv/bin/activate
   python3 -m pip install knita
   ```

4. Run your new pattern using the Knita CLI:
    ```bash
    knita build ./pattern.py
    ```

   You should see the following output:

    ```
   > knita build ./build/pattern.py
   example: finished
   ✓ hello-world (0s)
   
   Build log available at: /var/folders/bj/8x5csq0s15sgfpk2sq8ld7kw0000gn/T/knita/knita-build-20240517T034828Z.log
    ```

   And if you have a look at the build log:

   ```
   > cat /var/folders/bj/8x5csq0s15sgfpk2sq8ld7kw0000gn/T/knita/knita-build-20240517T034828Z.log                                                                                                                                                                                                                                                                                                                                          ─╯
   Pulling Docker image...
   Pulling image: docker.io/library/alpine:latest
   Using Docker registry auth: None
   Executing command: /bin/sh [-c echo "hello world"]
   hello world
   ```
   
