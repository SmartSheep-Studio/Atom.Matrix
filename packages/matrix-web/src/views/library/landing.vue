<template>
  <n-spin :show="reverting">
    <div class="container h-full">
      <div class="pt-10 pb-4 lg:px-10">
        <div class="text-xl font-bold">Library</div>
        <div>Browse items you own</div>
      </div>

      <n-grid :cols="4" item-responsive responsive="screen" x-gap="8" y-gap="8" class="lg:px-10 py-4">
        <n-gi span="4 m:2 l:1" v-for="item in library">
          <n-card class="cursor-pointer">
            <div class="text-lg">{{ apps[item.app_id].name }}</div>
            <div>{{ apps[item.app_id].description }}</div>
          </n-card>
        </n-gi>
      </n-grid>
    </div>
  </n-spin>
</template>

<script lang="ts" setup>
import { useMessage } from "naive-ui"
import { http } from "@/utils/http"
import { onMounted } from "vue"
import { ref } from "vue"

const $message = useMessage()

const reverting = ref(true)

const library = ref<any[]>([])
const apps = ref<{ [id: string]: any }>({})

async function fetch() {
  try {
    reverting.value = true
    library.value = (await http.get("/api/library")).data
    for (const item of library.value) {
      if (apps.value[item.app_id] == null) {
        apps.value[item.app_id] = (await http.get(`/api/explore/apps/${item.app_id}`)).data
      }
    }

    reverting.value = false
  } catch (e: any) {
    $message.error(`Something went wrong... ${e}`)
  }
}

onMounted(() => {
  fetch()
})
</script>
