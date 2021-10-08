
<div id="top"></div>

[![Go Reference](https://pkg.go.dev/badge/github.com/axpira/goplogjson.svg)](https://pkg.go.dev/github.com/axpira/goplogjson)
[![Go Report Card](https://goreportcard.com/badge/github.com/axpira/goplogjson)](https://goreportcard.com/report/github.com/axpira/goplogjson)
[![codecov](https://codecov.io/gh/axpira/goplogjson/branch/main/graph/badge.svg?token=FF2ZA1I0KX)](https://codecov.io/gh/axpira/goplogjson)
![Pipeline](https://github.com/axpira/goplogjson/actions/workflows/test.yml/badge.svg)

<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">

<h3 align="center">Gop Log Json</h3>
  <p align="center">
    A json implementation for [gop log](https://github.com/axpira/gop) with zero allocation
    <br />
    <a href="https://github.com/axpira/goplogjson"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/axpira/goplogjson/issues">Report Bug</a>
    ·
    <a href="https://github.com/axpira/goplogjson/issues">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li><a href="#installation">Installation</a></li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

A fast and with zero allocation implementation for [gop log](https://github.com/axpira/gop).

And no external depency, everything using builtin in Golang.

Output format is JSON and is very configurable

The main ideia is for [ZeroLog](https://github.com/rs/zerolog)

<p align="right">(<a href="#top">back to top</a>)</p>



### Built With

* [Golang](https://golang.org/)
* Zero allocation
* No external lib

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- INSTALLATION -->
## Installation

```sh
go get -u github.com/axpira/goplogjson
```

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

### Simple Log
```go
package main

import (
	"github.com/axpira/gop/log"
	_ "github.com/axpira/goplogjson"
)

func main() {
	log.Info("Hello World")
}

// Output: {"msg":"Hello World","level":"info","time":"2021-09-29T07:43:34-03:00"}
```

### Add fields to log
```go
package main

import (
	"github.com/axpira/gop/log"
	_ "github.com/axpira/goplogjson"
)

func main() {
	log.Inf(log.
		Str("str_field", "hello").
		Int("int_field", 42).
		Msg("Hello World"),
	)
}
// Output: {"str_field":"hello","int_field":42,"msg":"Hello World","level":"info","time":"2021-09-29T07:46:12-03:00"}
```

### Leveled logging
```go
package main

import (
	"errors"

	"github.com/axpira/gop/log"
	_ "github.com/axpira/goplogjson"
)

func main() {
	log.Error("Hello World", errors.New("unknown error"))
}
// Output: {"msg":"Hello World","err":"unknown error","level":"error","time":"2021-10-07T20:54:13-03:00"}
```

### Other Leveled logging
```go
package main

import (
	"errors"

	"github.com/axpira/gop/log"
	_ "github.com/axpira/goplogjson"
)

func main() {
	log.Log(
		log.ErrorLevel,
		log.Msg("Hello World").Err(errors.New("unknown error")),
	)

// Output: {"msg":"Hello World","err":"unknown error","level":"error","time":"2021-10-07T21:00:32-03:00"}
```

### Customize

```go
package main

import (
	"errors"
	"time"

	"github.com/axpira/gop/log"
	glj "github.com/axpira/goplogjson"
)

func main() {
	glj.DurationFieldUnit = time.Hour
	glj.ErrorFieldName = "myfielderr"
	glj.LevelFieldName = "mygreatlevel"
	glj.MessageFieldName = "mymessage"
	glj.TimeFormat = time.RFC1123
	glj.TimestampEnabled = false

	log.Log(
		log.ErrorLevel,
		log.Msg("Hello World").Err(errors.New("unknown error")),
	)
// Output: {"mymessage":"Hello World","myfielderr":"unknown error","mygreatlevel":"error"}
```

_For more examples like contextual log, please refer to the [Gop Log](https://github.com/axpira/gop)_

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Thiago Ferreira - thiagogbferreira@gmail.com

Project Link: [https://github.com/axpira/goplogjson](https://github.com/axpira/goplogjson)

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/axpira/goplogjson.svg?style=for-the-badge
[contributors-url]: https://github.com/axpira/goplogjson/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/axpira/goplogjson.svg?style=for-the-badge
[forks-url]: https://github.com/axpira/goplogjson/network/members
[stars-shield]: https://img.shields.io/github/stars/axpira/goplogjson.svg?style=for-the-badge
[stars-url]: https://github.com/axpira/goplogjson/stargazers
[issues-shield]: https://img.shields.io/github/issues/axpira/goplogjson.svg?style=for-the-badge
[issues-url]: https://github.com/axpira/goplogjson/issues
[license-shield]: https://img.shields.io/github/license/axpira/goplogjson.svg?style=for-the-badge
[license-url]: https://github.com/axpira/goplogjson/blob/main/LICENSE.txt
