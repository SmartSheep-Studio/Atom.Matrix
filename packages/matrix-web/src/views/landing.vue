<template>
  <div class="container">
    <div class="pt-12 px-10 flex items-center">
      <div>
        <img :src="Logo" width="80" height="80" />
      </div>
      <div class="mt-1 pl-4">
        <div class="text-4xl font-bold">Matrix</div>
        <div class="text-lg">Committed towards improving the application download and run experience.</div>
      </div>
    </div>

    <div class="pt-10 px-8">
      <div class="text-lg">Explore Queue</div>
      <n-grid
        :cols="4"
        x-gap="8"
        y-gap="8"
        item-responsive
        responsive="screen"
        class="py-3 cursor-pointer"
        style="margin: 0 -8px"
      >
        <n-gi span="4 m:2 l:1" v-for="item in exploreQueue">
          <n-card hoverable @click="$router.push({ name: 'store.info', params: { app: item.slug } })">
            <div class="text-lg">{{ item.name }}</div>
            <div class="line-clamp-2">{{ item.description }}</div>
          </n-card>
        </n-gi>
      </n-grid>
    </div>
  </div>
</template>

<script lang="ts" setup>
import Logo from "@/assets/icon.png"
import { http } from "@/utils/http"
import { useMessage } from "naive-ui"
import { onMounted, ref } from "vue"
import { useI18n } from "vue-i18n"

const { t } = useI18n()

const $message = useMessage()

const reverting = ref(true)

const exploreQueue = ref<any[]>([])

async function fetch() {
  try {
    reverting.value = true
    exploreQueue.value = (await http.get("/api/explore/apps")).data
  } catch (e: any) {
    $message.error(t('common.feedback.unknown-error', [e]))
  } finally {
    reverting.value = false
  }
}

onMounted(() => {
  fetch()
})
</script>
