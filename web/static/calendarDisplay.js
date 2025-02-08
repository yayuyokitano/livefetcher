/**
 * @typedef {Object} CalendarEvent
 * @property {string} id - The calendar id of the event
 * @property {string} name - The name of the event
 * @property {string} start - The time the event starts, as a string
 * @property {string} end - The time the event ends, as a string
 */

/**
 *
 * @param {string} openTime - The time the live opens, as a string
 * @param {string} startTime - The time the live starts, as a string
 * @param {Array.<string>} artists - An array of the artists performing at the live
 * @param {Object.<string, Array.<CalendarEvent>>} calendarEvents - The calendar events to look for conflicts with
 *
 * @returns {Array.<CalendarEvent>} - The list of conflicting events
 */
function getConflicts(openTime, startTime, artists, calendarEvents) {
  const boundaryStart = new Date(openTime);
  boundaryStart.setHours(boundaryStart.getHours() - 1);
  const boundaryEnd = new Date(startTime);
  boundaryEnd.setHours(
    boundaryEnd.getHours() + liveDurationHours(artists.length) + 1,
  );

  /**
   * @type {Object.<string, CalendarEvent>}
   */
  const conflictingEvents = {};
  for (
    let cur = new Date(boundaryStart.getTime());
    cur <= boundaryEnd;
    cur.setDate(cur.getDate + 1)
  ) {
    const curDayEvents =
      calendarEvents[
        `${cur.getFullYear()}-${(cur.getMonth() + 1).toString().padStart(2, "0")}-${cur.getDate().toString().padStart(2, "0")}`
      ];
    if (!curDayEvents) {
      continue;
    }
    outerLoop: for (const event of curDayEvents) {
      if (new Date(event.end) < boundaryStart) {
        continue;
      }
      if (new Date(event.start) > boundaryEnd) {
        continue;
      }
      if (event.name.startsWith("OPEN ")) {
        for (const e of Object.values(conflictingEvents)) {
          if (
            e.name.startsWith("START ") &&
            e.name.slice(6) === event.name.slice(5)
          ) {
            e.name = e.name.replace("START ", "");
            continue outerLoop;
          }
        }
      }
      if (event.name.startsWith("START ")) {
        for (const e of Object.values(conflictingEvents)) {
          if (
            e.name.startsWith("OPEN ") &&
            e.name.slice(5) === event.name.slice(6)
          ) {
            e.name = e.name.replace("OPEN ", "");
            continue outerLoop;
          }
        }
      }
      conflictingEvents[event.id] = { ...event };
    }
  }
  return Object.values(conflictingEvents);
}

/**
 *
 * @param {number} artistLength - The number of artists performing at the event
 * @returns An estimate of the number of hours the event lasts
 */
function liveDurationHours(artistLength) {
  switch (artistLength) {
    case 1:
      return 2;
    case 2:
      return 3;
    default:
      return Math.min(artistLength, 10);
  }
}
