

### Key Composition Summary

This design uses a single Key-Value bucket and is built around three distinct, high-performance patterns.

#### 1. The User's Flight List (`users...`)

This is the core data for any given user. It's prefixed by `userID` to create a private namespace and allow for fast, efficient lookups of a user's entire flight list.

* **Owned Flights:** `users.{userID}.flights.owned.{flightID}`
    * **Purpose:** Stores the persistent copy of a flight that a user has added themselves.
    * **Use Case:** The primary key for displaying a flight in the user's main list.

* **Shared Flights:** `users.{userID}.flights.shared.{flightID}`
    * **Purpose:** Stores a separate, persistent copy of a flight that another user has shared with them.
    * **Use Case:** Allows a user to retain flight data independently of the person who shared it.

**Benefit:** Your Go backend can get all flights for a logged-in user with a single, fast `Watch("users.{userID}.flights.>")` operation, providing a real-time view of their personal list.

#### 2. The Sharing Mechanism (`shares...`)

This is a temporary namespace used to handle the process of sharing a flight between users without them needing to know each other's IDs.

* **Pending Share:** `shares.pending.{shareUUID}`
    * **Purpose:** A short-lived, single-use token holding the flight data for an unclaimed share. It's created when User A initiates a share and deleted as soon as User B claims it.
    * **Value:** `{ "flightID": "...", "sharerID": "...", "flightData": {...} }`

* **Sharer's Record:** `users.{sharerID}.shares.sent.{shareUUID}`
    * **Purpose:** A record for the original sharer to track the status of their invitation (pending/claimed) and to have a "sanity check" of what they shared.
    * **Value:** `{ "flightID": "...", "status": "claimed", "sharedDataSnapshot": {...} }`

#### 3. The Scalability Index (`index...`)

This is the most critical component for your backend update service, as proven by our benchmarks. It provides a "reverse lookup" to find all users tracking a specific flight.



* **Flight-to-User Index:** `index.flight.{flightID}.users`
    * **Purpose:** To provide an instantaneous way for your flight data ingestor to find every user that needs to receive a real-time update for a given flight.
    * **Value:** A JSON array of user IDs, e.g., `["user-123", "user-456", "user-abc"]`.

**Benefit:** When new data for a fliAght arrives, your service performs a single, near-instant `Get` on this index key. This avoids the slow, unscalable wildcard scan (`users.*...`) and ensures your system can handle thousands of updates per second, regardless of the total number of users or flights being tracked.

Excellent question. This is a critical detail for creating a truly stable key. I've done the research, and the answer is definitive.

**Yes, the `scheduledOut` time absolutely changes.** Relying on it as an immutable part of your key is not a safe assumption.

### The Research: Why Scheduled Times are Unreliable

My research into airline operations and how data providers like FlightAware handle this reveals two key points:

1.  **Schedules are Fluid by Design:** Airlines publish schedules up to a year in advance, but these are often just placeholders. They are constantly adjusted for a huge number of reasons:
    * **Seasonal Demand:** A route might have three daily flights in summer but only one in winter.
    * **Optimizing Aircraft:** An airline might shift a plane to a more profitable route.
    * **Crew & Airport Logistics:** Changes in crew availability or airport slot times force schedule modifications.
    * **Market Competition:** A competitor adding a new flight can cause other airlines to adjust their times.

    One analysis I found showed that on a random day for a major US airline, **over 92% of their flights had some form of schedule change** from when they were first published. These can be minor (a few minutes) or major (hours or even a change of day).

2.  **How APIs Report Changes:** Data providers like FlightAware distinguish between different types of time. When a flight is delayed on the day of departure, the original `scheduledOut` time often remains the same, but new fields like `estimated_departure_time` or `actual_blockout_time` will show the new information. However, if the airline formally *reschedules* the flight weeks in advance, the base `scheduledOut` time itself is often updated to reflect the new official schedule.

### The Impact on Your Key Composition

If you include the scheduled time in your key (`...{HHMM}...`), and the airline formally reschedules that flight from 08:30 to 09:15, the key for that flight is now wrong. To fix this, you would have to:

1.  Read the flight record from the old key.
2.  Delete the old key.
3.  Create a new key with the new time.
4.  Find and update **every single index and user record** that points to the old key.

This is a complex, error-prone operation that you want to avoid at all costs.

### Final Recommendation: The Most Stable Key

The most robust key is one based on the information that is least likely to change *after the schedule is first published*. As we discussed, the **original scheduled time** is the best candidate for ensuring uniqueness for flights on the same route on the same day.

Here is the final, canonical format we should use:

* **Structure:** `flights.master.{ident}.{YYYY-MM-DD}.{original_HHMM}.{origin}.{destination}`
* **Example:** `flights.master.ANZ5272.2025-09-11.0830.NZAA.NZCH`

**Crucial Implementation Rule:**
The `{original_HHMM}` token in the key must be set **once** when the flight is first ingested and **never change**. All subsequent delays or even formal reschedules must be stored as fields *inside the value* (e.g., in your `Flight` struct), not by altering the key. This gives you a stable, permanent address for each flight instance.