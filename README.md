# synth [![CI](https://github.com/xyproto/synth/actions/workflows/ci.yml/badge.svg)](https://github.com/xyproto/synth/actions/workflows/ci.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/synth)](https://goreportcard.com/report/github.com/xyproto/synth)

Generate audio samples.

* This package is a bit experimental and a work in progress.
* Includes a `kick` utility for generating kick drum samples, a `snare` utility for generating snare drum samples, `rms` and `linear` for mixing audio and `sweep` for generating a detuned synth samples. They are all in the `cmd` directory.
* Used by the [Kickpad](https://github.com/xyproto/kickpad) application.
* Uses SDL2 to play sounds (unless `-tags ff` is given, then `ffplay` will be used instead).

## General info

* License: MIT
* Version: 1.13.0
