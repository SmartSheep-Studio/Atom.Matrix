<template>
  <div>
    <n-spin :show="requesting">
      <n-list bordered>
        <n-empty v-if="data.length <= 0" class="py-8" description="There's no data here. Why don't you create one?" />

        <n-list-item v-for="item in data">
          <n-thing :title="item.name">
            <template #header-extra>
              <div class="flex items-center">
                <div class="me-2">{{ item.is_published ? "Published" : "Draft" }}</div>
                <n-badge dot processing :type="item.is_published ? 'success' : 'warning'" />
              </div>
            </template>
            <template #description>
              <n-space size="small">
                <n-tag :bordered="false" type="success" size="small">#{{ item.slug }}</n-tag>
                <n-tag v-for="tag in item.post.tags" :bordered="false" type="primary" size="small">{{ tag }}</n-tag>
              </n-space>
            </template>

            <div>{{ item.description }}</div>
            <n-space class="mt-2" size="small">
              <n-button
                type="warning"
                size="small"
                @click="
                  $router.push({
                    name: 'console.apps.releases.update',
                    params: { app: props.data.slug, release: item.slug },
                  })
                "
              >
                <template #icon>
                  <n-icon :component="EditRound" />
                </template>
                Edit
              </n-button>
              <n-button type="error" size="small" @click="destroy(item)">
                <template #icon>
                  <n-icon :component="DeleteRound" />
                </template>
                Destroy
              </n-button>
            </n-space>
          </n-thing>
        </n-list-item>
      </n-list>
    </n-spin>

    <div class="flex justify-between mt-4">
      <n-button type="primary" size="small" @click="$router.push({ name: 'console.apps.releases.create' })">
        <template #icon>
          <n-icon :component="PlusRound" />
        </template>
        New Release
      </n-button>

      <n-pagination
        v-model:page="pagination.page"
        :page-count="Math.ceil((rawData.length ?? 0) / pagination.pageSize)"
        :page-slot="pagination.slot"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { usePrincipal } from "@/stores/principal"
import { computed, onMounted, reactive, ref } from "vue"
import { PlusRound, EditRound, DeleteRound } from "@vicons/material"
import { http } from "@/utils/http"
import { useDialog, useMessage } from "naive-ui"

const $dialog = useDialog()
const $message = useMessage()
const $principal = usePrincipal()

const props = defineProps<{ data: any }>()

const rawData = ref<any[]>([])
const data = computed(() => {
  const start = (pagination.page - 1) * pagination.pageSize
  return rawData.value?.reverse().slice(start, start + pagination.pageSize) ?? []
})

const requesting = ref(true)

const pagination = reactive({
  page: 1,
  pageSize: 5,
  slot: 5,
})

async function fetch() {
  try {
    requesting.value = true
    rawData.value = (await http.get(`/api/apps/${props.data.slug}/releases`)).data
  } catch (e: any) {
    $message.error(`Something went wrong... ${e}`)
  } finally {
    requesting.value = false
  }
}

function destroy(item: any) {
  $dialog.warning({
    title: "Warning",
    content: "This operation cannot be undo. Are you confirm?",
    positiveText: "Yes",
    negativeText: "Not really",
    onPositiveClick: async () => {
      try {
        requesting.value = true

        await http.delete(`/api/apps/${props.data.slug}/releases/${item.slug}`)
        await fetch()

        $message.success("Successfully deleted the post.")
      } catch (e: any) {
        $message.error(`Something went wrong... ${e}`)
      } finally {
        requesting.value = false
      }
    },
  })
}

onMounted(() => {
  fetch()
})
</script>
