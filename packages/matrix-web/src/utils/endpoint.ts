import { useEndpoint } from "@/stores/endpoint"

export function getEndpointPath(path: string, search?: string): string {
  const $endpoint = useEndpoint()

  return `${$endpoint.configuration.base_url}${path}${search ? "?" + search : ""}`
}
