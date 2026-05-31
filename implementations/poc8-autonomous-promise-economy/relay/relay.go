package relay

// AppName names the app-level relay role used between neighboring containers.
// Intent: Relay is an app promise to forward a message to a neighbor; it is not
// a kernel route command or global network authority. Source: DI-tugih
const AppName = "relay"
