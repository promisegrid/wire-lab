# SIM-gibut-conditional-release-group-session-local

This simulation turns the group-session-local conditional-release alternative
from `SIM-zarud-conditional-release-geofencing` into a concrete candidate
specimen. It tests whether onward-restraint, geofencing, and release conditions
belong inside the group/session semantics that deliver the content. Source:
`DI-fibuv`.

## Design Under Test

Group-session messages carry or reference release conditions, and each session
participant promises to preserve those conditions when dispatching content to
later recipients.

## Boundaries

This simulation does not modify the current group-session specimen. It tests
whether session-local ownership gives clear human semantics or overloads the
session layer with policy that should be separate.
