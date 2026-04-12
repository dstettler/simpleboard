const gameStatusMembers = ["InProgress", "Waiting", "Error"] as const;
export type GameStatus = typeof gameStatusMembers[number];

export function parseGameStatus(statusStr: string): GameStatus {
  const found = gameStatusMembers.find((matched) => matched === statusStr);
  if (found) {
    return found;
  }

  return "Error";
}
