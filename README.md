# voltorb_tilt
A voltorb flip solver for those who get tilted at the fact that this is picross that you can't perfect

## Things I've done to set this up:
```bash
go mod init github.com/megdorkable/voltorb_tilt
```

## Build
```bash
mkdir bin
```
```bash
go build -o bin/ .
```
Then run:
```bash
./bin/voltorb_tilt
```

## Run
```bash
go run .
```

## Flags
```bash
% ./bin/voltorb_tilt -h
Usage of ./bin/voltorb_tilt:
  -g    generate a random board
```