# Liar's Dice
A simple game of liar's dice for human, computer and AI.

## Game Rules
- Two or more players
- Each player starts with one or more dice each (generally five dice at start).
- The winner is the last player with at least one die.
- In each round,
  - Players take turns bidding how many dice of a specific die face there will be in sum (i.e. across all players).
  - Each subsequent player must bid a higher face value and/or a higher number of dice. For example if the current bid is 2 4s (i.e. in total across all players, there are __at least__ 2 dice with the value 4), then a player must bid higher face, e.g. 2 5s, or a higher number, e.g. 3 2s.
  - A player can, instead of raising, play two other actions which end the round
   - __Call__: Where the player does not believe the current bid is true (i.e. there are less than 2 4s).
   - __Exact__: Where the player believes the current bid is exactly true (i.e. there are exactly 2 4s, and not for example 3 4s).
 - On a call, if the player is wrong, they lose one die. If they are right, the previous player (who made the bet), loses one die.
 - On an exact, if the player is wrong, they lose one die. If they are right, they gain 2 dice. The previous player (who made the bet) is not effected.
 - Rounds are played, and people eliminated until only one person has a die/dice.

### Game Rule Translations
- [中文 Chinese](README-translated/README-Chinese.md) 🇨🇳

## Usage
First install
```shell
  go install github.com/Jeadie/liars-dice@latest
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
