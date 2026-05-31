package issuer

// AppName is the script-visible issuer app role for local token promises.
// Intent: Keep issuing behavior explicitly app-level so tokens remain promises
// by resource-controlling agents instead of kernel-granted grants. Source:
// DI-tugih; DI-tanat
const AppName = "issuer"
