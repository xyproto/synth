# synth

Generate audio samples.

* This package is a bit experimental and a work in progress.
* Includes a `kick` utility for generating kick drum samples, plus `rms` and `linear` for mixing audio and `sweep` for making a detuned synth sample. They are all in the `cmd` directory.
* Used by the [Kickpad](https://github.com/xyproto/kickpad) application.
* Uses SDL2 to play sounds (unless `-tags ff` is given, then `ffplay` will be used instead).

## General info

* License: MIT
* Version: 1.12.2
