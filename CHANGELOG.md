# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



## [Unreleased]

- Add tests for quoted URL matching.

## [0.3.0] - 2021-02-18

### Added

- Add new error regex for internal API errors often returned when an API is not yet ready to handle requests.



## [0.2.3] 2020-04-17

### Fixed

- Fix error message prefix matching in IsAPINotAvailable().



## [0.2.2] 2020-04-06

### Fixed

- Fix all error matchings for changing error messages by considering optional quotes.


## [0.2.1] 2020-04-03

### Fixed

- Fix error matching for changing error messages.



## [0.2.0] 2020-03-25

### Changed

- Switch from dep to go modules
- Use architect-orb



## [0.1.0] 2020-03-19

### Added

- First release.



[Unreleased]: https://github.com/giantswarm/errors/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/giantswarm/errors/compare/v0.2.3...v0.3.0
[0.2.3]: https://github.com/giantswarm/errors/compare/v0.2.2...v0.2.3
[0.2.2]: https://github.com/giantswarm/errors/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/giantswarm/errors/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/errors/compare/v0.1.0...v0.2.0

[0.1.0]: https://github.com/giantswarm/errors/releases/tag/v0.1.0
