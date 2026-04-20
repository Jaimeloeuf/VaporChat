export function isIsoDatetimeOlderThan(
  isoDatetime: string,
  olderThanTimeInSeconds: number,
): boolean {
  const inputTime = new Date(isoDatetime).getTime()
  if (isNaN(inputTime)) {
    throw new Error('Invalid ISO datetime string')
  }

  const currentTime = Date.now()
  return currentTime - inputTime > olderThanTimeInSeconds
}
