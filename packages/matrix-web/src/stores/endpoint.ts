import { defineStore } from "pinia"
import { http } from "@/utils/http"
import { ref } from "vue"

export const useEndpoint = defineStore("endpoint", () => {
  const isPrepared = ref(true)
  const service = ref<any>({})
  const additional = ref<any>({})
  const configuration = ref<any>({})

  async function fetch() {
    try {
      const res = await http.get("/api/info")
      service.value = res.data.service
      additional.value = res.data.additional
      configuration.value = res.data.configuration
    } catch (e: any) {
      if (e.response.status === 503) {
        isPrepared.value = false
      } else {
        isPrepared.value = true
      }
      throw e
    }

    document.title = configuration.value.general.name ?? "Project Atom"
  }

  return { isPrepared, fetch }
})
