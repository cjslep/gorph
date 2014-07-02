Table of Contents
---

1. [About](#about)
2. [Dependencies](#dependencies)
3. [License](#license)
4. [Installation](#installation)
5. [Usage/API](#usageapi)
6. [Contributing](#contributing)
7. [Authors](#authors)

<a name="about"/>
About
---

Gorph is a library that manipulates graphic images in order to procedurally modify or create derived images.

It is currently under development with the goal of enabling server-side clients to pre-render or procedurally generate images.

The name is an homage to the library itself. It is derived from using the "morph" of two images to create a keyframe-based animation system. It also comes from "gopher" or "golang". *Morphing* the two different words gives the perhaps-not-quite-as-humorous result.

This library is still unstable, please do not hesitate to contact me with questions. Comments in the source are nonexistant, and is the first item to be remedied.

<a name="license"/>
License
---

Gorph is released under the [MIT Expat License](./LICENSE), a permissive free software license. Contributing back to the source is appreciated greatly, but not dictated by the license. A copy of the license is distributed with the source code.

<a name="dependencies"/>
Dependencies
---

At this time, no additional dependencies are required beyond the golang standard library.

<a name="installation"/>
Installation
---

If Go is not installed, (please do so first)[http://golang.org/doc/install].

To download the source, use the command:
```
go get github.com/cjslep/gorph
```

Then, to build the source:
```
go build github.com/cjslep/gorph
```

To ensure all tests pass, run:
```
go test github.com/cjslep/gorph
```

Note that the `dev-unstable` branch may not pass all tests.

<a name="usageapi"/>
Usage / API
---

This section will be more detailed once development stabilizes.

<a name="contributing"/>
Contributing
---

The license governing this source code does not require users to contribute back to this source code, nor does it prevent users from creating proprietary derivative works. However, contributions back to the source are greatly appreciated.

Comments in the source are nonexistant, and is the first item to be remedied.

<a name="authors"/>
Authors
---

In no particular order:
* [cjslep](https://github.com/cjslep)