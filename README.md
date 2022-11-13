# Liar's Dice
A simple game of liar's dice for human, computer and AI.

## Usage
First install
```shell
  go install github.com/Jeadie/liars-dice
```
To play a simple game, two users starting with three dice each, against algorithmic opponents:
```shell
   liars-dice --user-idx 0 -- 3 3
```
Or just a single round
```shell
  liars-dice --user-idx 0 round -- 3 3
```

Or setup a game over the network, start the server
```shell
  liars-dice --ws-agents 1  -- 3 3
```
And connect on the other side
```shell
  liars-dice  ws-client
```
