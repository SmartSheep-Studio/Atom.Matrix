import { type RouteLocationNormalizedLoaded } from "vue-router"
import { usePrincipal } from "@/stores/principal"
import wildcard from "wildcard-match"

export function hasUserPermissions(...requires: string[]) {
  const $principal = usePrincipal()
  if (!$principal.isLoggedIn || $principal.account == null) {
    return false
  }

  for (const require of requires) {
    let passed = false
    for (const perm of $principal.account.permissions ?? []) {
      if (wildcard(perm)(require)) {
        passed = true
        break
      }
    }

    if (!passed) {
      return false
    }
  }

  return true
}

export function keepGate(to: RouteLocationNormalizedLoaded) {
  const $principal = usePrincipal()
  const meta: any = to?.meta?.gatekeeper ?? {}

  if (meta?.must === true && !$principal.isLoggedIn) {
    return false
  } else if (meta?.permissions != null && !hasUserPermissions(...meta?.permissions)) {
    return false
  }

  return true
}
