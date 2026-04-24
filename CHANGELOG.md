# Changelog

## [0.2.0-alpha.1](https://github.com/michaeldcanady/go-onedrive/compare/v0.2.0-alpha...v0.2.0-alpha.1) (2026-04-24)


### Features

* **fs:** uri factory ([#122](https://github.com/michaeldcanady/go-onedrive/issues/122)) ([198b750](https://github.com/michaeldcanady/go-onedrive/commit/198b7504440db332a96827331134df9d5e1d226b))
* redesign drive management and access design ([#126](https://github.com/michaeldcanady/go-onedrive/issues/126)) ([a1cf13d](https://github.com/michaeldcanady/go-onedrive/commit/a1cf13dd6252393ae42077d2622afc600c42b7f9))
* sort spec ([#121](https://github.com/michaeldcanady/go-onedrive/issues/121)) ([e966c5e](https://github.com/michaeldcanady/go-onedrive/commit/e966c5e4f7661d0ddb78ddd306f32d1ee959e845))
* **spec:** add filter spec ([#120](https://github.com/michaeldcanady/go-onedrive/issues/120)) ([fa7ce0e](https://github.com/michaeldcanady/go-onedrive/commit/fa7ce0ec152d51dcfb0c97bea3d7af84d2ad44e3))

## [0.2.0-alpha](https://github.com/michaeldcanady/go-onedrive/compare/v0.1.0...v0.2.0-alpha) (2026-04-07)


### Features

* add support for env vars ([#99](https://github.com/michaeldcanady/go-onedrive/issues/99)) ([caaf55a](https://github.com/michaeldcanady/go-onedrive/commit/caaf55a20fbfce4b9b9d73329b7ec6409aec83d5))
* **cli:** add mv command for moving and renaming items ([#87](https://github.com/michaeldcanady/go-onedrive/issues/87)) ([46d37cd](https://github.com/michaeldcanady/go-onedrive/commit/46d37cd867fbbf3044603d766d3410a0d3ab2d70))
* **cli:** add touch command for empty file creation ([#86](https://github.com/michaeldcanady/go-onedrive/issues/86)) ([286e586](https://github.com/michaeldcanady/go-onedrive/commit/286e5863cf2c8c0d931614001c6eda3228ebbf4b))
* **cli:** handle SIGINT with ErrUserAbort in main ([d7098ee](https://github.com/michaeldcanady/go-onedrive/commit/d7098eec29c5707e9f7464007a0046447931981b))
* config commands ([#110](https://github.com/michaeldcanady/go-onedrive/issues/110)) ([dd4ccff](https://github.com/michaeldcanady/go-onedrive/commit/dd4ccff35ebc69f7511ae7170410a9ef4b8006bb))
* config profile unification ([#108](https://github.com/michaeldcanady/go-onedrive/issues/108)) ([13ceba6](https://github.com/michaeldcanady/go-onedrive/commit/13ceba6c9924f3fd1756170f6b9f3802b5275ed5))
* onedrive vertical slice ([#98](https://github.com/michaeldcanady/go-onedrive/issues/98)) ([d7098ee](https://github.com/michaeldcanady/go-onedrive/commit/d7098eec29c5707e9f7464007a0046447931981b))
* **rm:** add rm command ([#93](https://github.com/michaeldcanady/go-onedrive/issues/93)) ([053d397](https://github.com/michaeldcanady/go-onedrive/commit/053d397436fe947a44411b14eb9ba8a0cdd8d4e4))
* test prerelease ([f3bb1d5](https://github.com/michaeldcanady/go-onedrive/commit/f3bb1d52fbb9695b62d168c50fb19be9356b8601))


### Bug Fixes

* **config:** fill-in missing configs with defaults ([1f85df7](https://github.com/michaeldcanady/go-onedrive/commit/1f85df781fe9b2eaa73a09a99c4c332842299fbf))
* fix ls only returns os args ([c774592](https://github.com/michaeldcanady/go-onedrive/commit/c7745929c3df7a356512540152213d5f26983096))

## [0.1.1](https://github.com/michaeldcanady/go-onedrive/compare/v0.1.0...v0.1.1) (2026-04-05)


### Features

* Add support for env vars ([#99](https://github.com/michaeldcanady/go-onedrive/issues/99)) ([caaf55a](https://github.com/michaeldcanady/go-onedrive/commit/caaf55a20fbfce4b9b9d73329b7ec6409aec83d5))
* **cli:** Add mv command for moving and renaming items ([#87](https://github.com/michaeldcanady/go-onedrive/issues/87)) ([46d37cd](https://github.com/michaeldcanady/go-onedrive/commit/46d37cd867fbbf3044603d766d3410a0d3ab2d70))
* **cli:** Add touch command for empty file creation ([#86](https://github.com/michaeldcanady/go-onedrive/issues/86)) ([286e586](https://github.com/michaeldcanady/go-onedrive/commit/286e5863cf2c8c0d931614001c6eda3228ebbf4b))
* **cli:** Handle SIGINT with ErrUserAbort in main ([d7098ee](https://github.com/michaeldcanady/go-onedrive/commit/d7098eec29c5707e9f7464007a0046447931981b))
* Config commands ([#110](https://github.com/michaeldcanady/go-onedrive/issues/110)) ([dd4ccff](https://github.com/michaeldcanady/go-onedrive/commit/dd4ccff35ebc69f7511ae7170410a9ef4b8006bb))
* Config profile unification ([#108](https://github.com/michaeldcanady/go-onedrive/issues/108)) ([13ceba6](https://github.com/michaeldcanady/go-onedrive/commit/13ceba6c9924f3fd1756170f6b9f3802b5275ed5))
* Onedrive vertical slice ([#98](https://github.com/michaeldcanady/go-onedrive/issues/98)) ([d7098ee](https://github.com/michaeldcanady/go-onedrive/commit/d7098eec29c5707e9f7464007a0046447931981b))
* **rm:** Add rm command ([#93](https://github.com/michaeldcanady/go-onedrive/issues/93)) ([053d397](https://github.com/michaeldcanady/go-onedrive/commit/053d397436fe947a44411b14eb9ba8a0cdd8d4e4))


### Bug Fixes

* **config:** Fill-in missing configs with defaults ([1f85df7](https://github.com/michaeldcanady/go-onedrive/commit/1f85df781fe9b2eaa73a09a99c4c332842299fbf))
* Fix ls only returns os args ([c774592](https://github.com/michaeldcanady/go-onedrive/commit/c7745929c3df7a356512540152213d5f26983096))
