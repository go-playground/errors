# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [5.3.3] - 2023-08-16
### Fixed
- First error inconsistently wrapped with error and prefix instead of err then prefix in the chain.
- Corrected Is & As functions to directly causal error.

## [5.3.2] - 2023-08-16
### Fixed
- Link Error output missing error when prefix was blank.

## [5.3.1] - 2023-08-15
### Fixed
- Wrap recursively wrapping the Chain itself instead of only adding another Link.

## [5.3.0] - 2023-04-14
### Fixed
- Added Error interface for Link.

## [5.2.3] - 2022-05-30
### Fixed
- Fixed errors.As wrapper, linter is wrong.

## [5.2.2] - 2022-05-30
### Fixed
- Fixed calling error helpers on first wrap only, where the source of the data will be.

## [5.2.1] - 2022-05-30
### Fixed
- Changelog pointing to wrong repo.

## [5.2.0] - 2022-05-28
### Added
- Wrap check for nil errors which will now panic to indicate a larger issue, the caller NOT checking an error value. This caused a panic when trying to print the error if nil error was wrapped.

### Removed
- Deprecated information in the documentation.

### Changed
- Updated deps.


[Unreleased]: https://github.com/go-playground/errors/compare/v5.3.3...HEAD
[5.3.3]: https://github.com/go-playground/errors/compare/v5.3.2...v5.3.3
[5.3.2]: https://github.com/go-playground/errors/compare/v5.3.1...v5.3.2
[5.3.1]: https://github.com/go-playground/errors/compare/v5.3.0...v5.3.1
[5.3.0]: https://github.com/go-playground/errors/compare/v5.2.3...v5.3.0
[5.2.3]: https://github.com/go-playground/errors/compare/v5.2.2...v5.2.3
[5.2.2]: https://github.com/go-playground/errors/compare/v5.2.1...v5.2.2
[5.2.1]: https://github.com/go-playground/errors/compare/v5.2.0...v5.2.1
[5.2.0]: https://github.com/go-playground/errors/compare/v5.1.1...v5.2.0
