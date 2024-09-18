#!/bin/bash

mkdir -p kicks
echo "Generating a variety of kick drum samples..."

# Generate classic kicks with adjustments

go run ../cmd/kick/main.go --808 -o kicks/808_classic.wav --length 650 --quality 96 --waveform 0 --attack 0.02 --decay 0.7 --sustain 0.15 --release 0.4 --sweep 0.95 --filter 4500 --pitchdecay 0.6 --drive 0.25 --bitdepth 16 --numoscillators 2 --oscillatorlevels 1.0,0.7 --saturator 0.3
go run ../cmd/kick/main.go --909 -o kicks/909_classic.wav --length 700 --quality 96 --waveform 1 --attack 0.002 --decay 0.55 --sustain 0.1 --release 0.35 --sweep 0.75 --filter 9000 --pitchdecay 0.35 --drive 0.45 --bitdepth 16 --numoscillators 2 --oscillatorlevels 1.0,0.6 --saturator 0.35
go run ../cmd/kick/main.go --707 -o kicks/707_classic.wav --length 600 --quality 96 --waveform 1 --attack 0.005 --decay 0.4 --sustain 0.25 --release 0.35 --sweep 0.65 --filter 5500 --pitchdecay 0.45 --drive 0.35 --bitdepth 16 --numoscillators 1 --oscillatorlevels 1.0 --saturator 0.25
go run ../cmd/kick/main.go --606 -o kicks/606_classic.wav --length 700 --quality 96 --waveform 0 --attack 0.01 --decay 0.35 --sustain 0.1 --release 0.25 --sweep 0.75 --filter 5500 --pitchdecay 0.55 --drive 0.45 --bitdepth 16 --numoscillators 2 --oscillatorlevels 1.0,0.85 --saturator 0.4

# Generate experimental samples

go run ../cmd/kick/main.go --deephouse -o kicks/deephouse_classic.wav --length 700 --quality 96 --waveform 0 --attack 0.007 --decay 0.9 --sustain 0.35 --release 0.75 --sweep 0.85 --filter 3700 --pitchdecay 0.65 --drive 0.65 --bitdepth 16 --numoscillators 3 --oscillatorlevels 1.0,0.75,0.5 --saturator 0.55
go run ../cmd/kick/main.go --experimental -o kicks/experimental_classic.wav --length 750 --quality 96 --waveform 2 --attack 0.001 --decay 0.8 --sustain 0.05 --release 0.5 --sweep 1.3 --filter 2800 --pitchdecay 0.9 --drive 0.9 --bitdepth 16 --numoscillators 4 --oscillatorlevels 1.0,0.5,0.3,0.2 --saturator 0.75 --filterbands 500,2000,4000

echo "Generating varied kick drum samples..."

# 10 varied kicks demonstrating different parameter combinations

go run ../cmd/kick/main.go --808 -o kicks/808_varied_1.wav --length 650 --quality 96 --waveform 0 --attack 0.015 --decay 0.65 --sustain 0.25 --release 0.55 --sweep 0.8 --filter 4700 --pitchdecay 0.55 --drive 0.35 --bitdepth 16 --numoscillators 2 --oscillatorlevels 1.0,0.75 --saturator 0.35
go run ../cmd/kick/main.go --909 -o kicks/909_varied_1.wav --length 700 --quality 96 --waveform 1 --attack 0.0015 --decay 0.45 --sustain 0.1 --release 0.25 --sweep 0.85 --filter 6500 --pitchdecay 0.4 --drive 0.55 --bitdepth 16 --numoscillators 3 --oscillatorlevels 1.0,0.55,0.3 --saturator 0.45
go run ../cmd/kick/main.go --606 -o kicks/606_varied_1.wav --length 700 --quality 96 --waveform 0 --attack 0.02 --decay 0.45 --sustain 0.2 --release 0.4 --sweep 0.85 --filter 5100 --pitchdecay 0.65 --drive 0.55 --bitdepth 16 --numoscillators 1 --oscillatorlevels 1.0 --saturator 0.25
go run ../cmd/kick/main.go --deephouse -o kicks/deephouse_varied_1.wav --length 900 --quality 96 --waveform 0 --attack 0.008 --decay 0.75 --sustain 0.3 --release 0.85 --sweep 0.95 --filter 3200 --pitchdecay 0.75 --drive 0.75 --bitdepth 16 --numoscillators 2 --oscillatorlevels 1.0,0.8 --saturator 0.55
go run ../cmd/kick/main.go --experimental -o kicks/experimental_varied_1.wav --length 900 --quality 96 --waveform 3 --attack 0.006 --decay 0.75 --sustain 0.1 --release 0.45 --sweep 1.1 --filter 2900 --pitchdecay 0.85 --drive 0.8 --bitdepth 16 --numoscillators 4 --oscillatorlevels 1.0,0.5,0.4,0.2 --saturator 0.8 --filterbands 400,1800,3500
go run ../cmd/kick/main.go --707 -o kicks/707_varied_1.wav --length 800 --quality 96 --waveform 1 --attack 0.004 --decay 0.6 --sustain 0.2 --release 0.4 --sweep 0.7 --filter 5100 --pitchdecay 0.35 --drive 0.55 --bitdepth 16 --numoscillators 2 --oscillatorlevels 1.0,0.65 --saturator 0.35
go run ../cmd/kick/main.go --deephouse -o kicks/deephouse_varied_2.wav --length 900 --quality 96 --waveform 0 --attack 0.01 --decay 0.9 --sustain 0.4 --release 0.9 --sweep 0.95 --filter 3500 --pitchdecay 0.65 --drive 0.85 --bitdepth 16 --numoscillators 2 --oscillatorlevels 1.0,0.9 --saturator 0.55

# Video game-style kicks

go run ../cmd/kick/main.go --experimental -o kicks/videogame_kick_1.wav --length 400 --quality 44 --waveform 2 --attack 0.001 --decay 0.1 --sustain 0.0 --release 0.05 --sweep 2.0 --filter 2500 --pitchdecay 1.3 --drive 0.5 --bitdepth 16 --numoscillators 1 --oscillatorlevels 1.0 --saturator 0.3
go run ../cmd/kick/main.go --experimental -o kicks/videogame_kick_2.wav --length 400 --quality 44 --waveform 3 --attack 0.001 --decay 0.12 --sustain 0.0 --release 0.06 --sweep 2.5 --filter 1800 --pitchdecay 1.7 --drive 0.6 --bitdepth 16 --numoscillators 1 --oscillatorlevels 1.0 --saturator 0.35

ffmpeg || exit 1

# Convert all .wav files to .mp4 using ffmpeg

for wav_file in kicks/*.wav; do
  mp4_file="${wav_file%.wav}.mp4"
  echo "Converting $wav_file to $mp4_file..."
  ffmpeg -i "$wav_file" -c:a aac "$mp4_file" -y
done

echo "All kick samples generated, converted to .mp4, and stored in the 'kicks' directory."
