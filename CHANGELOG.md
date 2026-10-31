# Changelog

## [2.0.0](https://github.com/eikendev/hackenv/compare/v1.0.0...v2.0.0) (2026-10-31)


### ⚠ BREAKING CHANGES

* **up:** drop the display size setting
* **get:** import local images unverified

### Features

* **get:** download the tuned parrot image ([8d45db3](https://github.com/eikendev/hackenv/commit/8d45db300130d71d2e51da7d1c7bb915a751657d))
* **get:** import local images unverified ([6e3556f](https://github.com/eikendev/hackenv/commit/6e3556f7edcd41709c59bf71b00074ea664ce2d4))
* **up:** drop the display size setting ([d732998](https://github.com/eikendev/hackenv/commit/d732998279934a9df58baddbfc60d730693449e0))


### Bug Fixes

* **images:** drop warning for uneven versions ([b543bef](https://github.com/eikendev/hackenv/commit/b543bef83ecb66cd056142b8c1d81e874b9017fa))
* **up:** provision parrot 7 live sessions ([4c02989](https://github.com/eikendev/hackenv/commit/4c029893cd3da1e64385f20bae1b5142ad8c936c))


### Documentation

* restructure readme around a quick start ([9eddfbb](https://github.com/eikendev/hackenv/commit/9eddfbb7f2e473a629cc0da86a0cc4d022971dd1))


### Build

* **actions:** bump zizmorcore/zizmor-action ([#37](https://github.com/eikendev/hackenv/issues/37)) ([4b5d40b](https://github.com/eikendev/hackenv/commit/4b5d40bf189a0dfe2e23e1ff45c0d5dc95219af4))
* **go:** bump the go-dependencies group across 1 directory with 3 updates ([#38](https://github.com/eikendev/hackenv/issues/38)) ([d3a3c5f](https://github.com/eikendev/hackenv/commit/d3a3c5fb61988ffabb0deb5fda5ae167dc205e4f))

## [1.0.0](https://github.com/eikendev/hackenv/compare/v0.3.0...v1.0.0) (2026-10-30)


### ⚠ BREAKING CHANGES

* drop guest keyboard layout setup
* accept memory units, rename display flag
* modernize vm definition
* use passt networking instead of a bridge

### Features

* accept memory units, rename display flag ([e37c420](https://github.com/eikendev/hackenv/commit/e37c4205d3119e316d2261e134d811f22d413d35))
* add doctor command ([5a81be0](https://github.com/eikendev/hackenv/commit/5a81be0e1e5a37dd50e9d3d4bee01f202f9da146))
* cancel downloads on interrupt ([5e7a1b6](https://github.com/eikendev/hackenv/commit/5e7a1b6fcc64e0ddc91d56855f153e694d2c1502))
* **doctor:** warn before signing keys expire ([21b3889](https://github.com/eikendev/hackenv/commit/21b38899bff9819325cd17cd8d7ca9cbe777e606))
* download kali images via torrent ([02aa7d1](https://github.com/eikendev/hackenv/commit/02aa7d1c8b61ad9732251733be6885fad4d899f6))
* drop guest keyboard layout setup ([fe4f787](https://github.com/eikendev/hackenv/commit/fe4f787de4b583e969c3a57bf24406254dd26b16))
* fall back when setxkbmap is missing ([45a0ba8](https://github.com/eikendev/hackenv/commit/45a0ba83df1749e9ddd3fed761e34b166556322b))
* import images with get --file ([a4ff0d8](https://github.com/eikendev/hackenv/commit/a4ff0d878dc9dbb82d605531b9bffd704b90a4b3))
* load options from config file ([f6d23c2](https://github.com/eikendev/hackenv/commit/f6d23c21b62bad867e5b42316c0ac9be0121a54e))
* modernize vm definition ([bdad49a](https://github.com/eikendev/hackenv/commit/bdad49a205d031b6bd89045326faed89ab01f379))
* use passt networking instead of a bridge ([39a6264](https://github.com/eikendev/hackenv/commit/39a62640d785d61934d578ff0ac169183e73be5f))
* verify image checksum signatures ([76c09b7](https://github.com/eikendev/hackenv/commit/76c09b763eddc20390f858642e350fe708070189))


### Bug Fixes

* **ci:** install libvirt headers for builds ([0a01ced](https://github.com/eikendev/hackenv/commit/0a01cedd615a0956ef5fb25ab56950cebe179530))
* **deps:** bump vulnerable torrent dependencies ([32ded0b](https://github.com/eikendev/hackenv/commit/32ded0be1782cc207eaca1bbfee7667218f8deb4))
* **get:** finish torrent when no bytes are missing ([c26c36b](https://github.com/eikendev/hackenv/commit/c26c36b437448f66dea0cf0e72f29a8b388b1a40))
* **get:** locate and resume torrent downloads ([b57cbc2](https://github.com/eikendev/hackenv/commit/b57cbc28d2fd8dfd35ee20705a8ae95b092cc228))
* **images:** discover latest parrot release ([eab6127](https://github.com/eikendev/hackenv/commit/eab6127ecd797444ebcd3af801369378e1d89999))
* **scripts:** skip labels without selinux ([47cb7ec](https://github.com/eikendev/hackenv/commit/47cb7ecd21a924701c088893e220e2f99e02dc28))
* stream script output to terminal ([6daf458](https://github.com/eikendev/hackenv/commit/6daf45803efc4ad32a33475afd4a22cde6d4a9b0))
* **up:** ignore eof when closing ssh sessions ([01e19e1](https://github.com/eikendev/hackenv/commit/01e19e15a50232a8256e37c5c17347524932eb51))
* **up:** make guest setup commands idempotent ([ba5aecb](https://github.com/eikendev/hackenv/commit/ba5aecb3999be4b4a8d0b62dadd06fe1dd9cc8fe))
* **up:** map guest user onto host for virtiofs ([e582abd](https://github.com/eikendev/hackenv/commit/e582abdd4fe00e27d1ae998a84038596da193929))
* **up:** require memory units and cap at host ram ([e321f40](https://github.com/eikendev/hackenv/commit/e321f40cba15f0aaaa8dbd01486109ae6873000e))
* **up:** require ssh banner before configuring ([ee118ce](https://github.com/eikendev/hackenv/commit/ee118ce7efa4c0d893a1f5d46c5cf2edad4a7d72))
* **up:** share directory via 9p only ([e3284e4](https://github.com/eikendev/hackenv/commit/e3284e4cdbe63fe8249817142ac9dd80cd97b209))


### Documentation

* describe supported distributions ([ecc18ab](https://github.com/eikendev/hackenv/commit/ecc18ab8dd1dada637ea4e31410b6e16aa9cf8e4))
* refresh installation and sharing notes ([78eac0b](https://github.com/eikendev/hackenv/commit/78eac0b06f2cd2748de99536496b3aad5e15af77))


### Build

* **actions:** bump zizmorcore/zizmor-action ([#34](https://github.com/eikendev/hackenv/issues/34)) ([3b83ea4](https://github.com/eikendev/hackenv/commit/3b83ea480cdf15cf2b03fc3db3fdcfd2c84bfba4))
* **go:** bump go.opentelemetry.io/otel from 1.38.0 to 1.41.0 ([#36](https://github.com/eikendev/hackenv/issues/36)) ([f5f3a47](https://github.com/eikendev/hackenv/commit/f5f3a473af49e8e24016f0e3965dcc0dd19b0e8e))
* load libvirt at runtime via dlopen ([1a88a65](https://github.com/eikendev/hackenv/commit/1a88a6506688e309f3d00a34680985e37eed4f8d))
* pin lint tools in a tools module ([e882619](https://github.com/eikendev/hackenv/commit/e8826195fc5082c2078045a71ae55fcf11ec7d76))
* release linux arm64 binaries ([dc47d81](https://github.com/eikendev/hackenv/commit/dc47d81901eff6a09f3deac426b3a7ff2fe456b9))
* use pure go torrent backends ([9116418](https://github.com/eikendev/hackenv/commit/9116418fe7668246a780de6f5e3899ae57abe7d0))

## [0.3.0](https://github.com/eikendev/hackenv/compare/v0.2.3...v0.3.0) (2026-09-19)


### Features

* attest build provenance of release artifacts ([112c15d](https://github.com/eikendev/hackenv/commit/112c15d01f4b00a9e7bc7b7014efb08de2ee316a))


### Bug Fixes

* **deps:** bump go to 1.27 and x/crypto to v0.56.0 ([9035bea](https://github.com/eikendev/hackenv/commit/9035beafb848d8b13eed4139a42865781eb52a73))


### Build

* route dependabot commits to the build type ([e68e359](https://github.com/eikendev/hackenv/commit/e68e359d3b7bfcd6443dcabdef7e9b2a8bad09d3))

## [0.2.3](https://github.com/eikendev/hackenv/compare/v0.2.2...v0.2.3) (2026-03-21)


### Bug Fixes

* avoid ipv4-mapped ipv6 output ([ceb326f](https://github.com/eikendev/hackenv/commit/ceb326fb0289bede292a1b95388378c46b31f8c5))
* handle mismatched version parts safely ([f25b76a](https://github.com/eikendev/hackenv/commit/f25b76ab89d0df26cfd7ceb926883193d8ce9ba3))
* honor strictness when fetching image metadata ([7e540ec](https://github.com/eikendev/hackenv/commit/7e540ec547ef3a4b0b1d1e24e2f644f888ebd023))
* **libvirt:** handle missing interface addresses ([406a2b3](https://github.com/eikendev/hackenv/commit/406a2b3eb228cd978f7c08ac5f547bfebc41317a))


### Documentation

* fix typo in GetDownloadInfo comment ([0133b5d](https://github.com/eikendev/hackenv/commit/0133b5d014ab80e5c56d6d6a85a356feab5aa001))


### Build

* enforce go fix check in test ([4548f60](https://github.com/eikendev/hackenv/commit/4548f60b67517c09838f7a49a9e13a65ede0a55b))
* include docs in release archive ([202ff7c](https://github.com/eikendev/hackenv/commit/202ff7c9df1d4f39fce02293c4f1c78c0e42b352))

## [0.2.2](https://github.com/eikendev/hackenv/compare/v0.2.1...v0.2.2) (2026-02-16)


### Bug Fixes

* trigger build with clean version command ([db519b2](https://github.com/eikendev/hackenv/commit/db519b2d3e71ff5eb3bd6cbe4d7670eac5bf5af4))


### Documentation

* clean up changelog ([6fe43b7](https://github.com/eikendev/hackenv/commit/6fe43b75a15ab6e2757e2bec196cac3d2b98253c))
* update build badge link ([47c2315](https://github.com/eikendev/hackenv/commit/47c2315c068d828db1a698fe6c849d32357ee4bd))

## [0.2.1](https://github.com/eikendev/hackenv/compare/v0.2.0...v0.2.1) (2026-02-16)


### Documentation

* initialize changelog ([060c75e](https://github.com/eikendev/hackenv/commit/060c75e889c923a76d6c0bb1799664afe7c8a552))
