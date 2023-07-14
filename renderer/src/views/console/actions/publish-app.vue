<template>
  <n-card embedded>
    <div class="text-lg">{{ keyword }} App</div>
    <div class="mt-1">
      Do you want {{ keyword }} this application? This operation will change the content did users view.
    </div>

    <div class="mt-3">
      <n-button
        class="w-full"
        :type="props.data.is_published ? 'warning' : 'primary'"
        :loading="submitting"
        @click="publish"
      >
        <template #icon>
          <n-icon :component="keyicon" />
        </template>
        {{ keyword }}
      </n-button>
    </div>
  </n-card>
</template>

<script lang="ts" setup>
import { http } from "@/utils/http"
import { CloudUploadRound, CloudDownloadRound } from "@vicons/material"
import { useDialog, useMessage } from "naive-ui"
import { computed, ref } from "vue"

const $dialog = useDialog()
const $message = useMessage()

const props = defineProps<{ data: any }>()
const emits = defineEmits(["refresh"])

const submitting = ref(false)

const keyword = computed(() => (props.data.is_published ? "Unpublish" : "Publish"))
const keyicon = computed(() => (props.data.is_published ? CloudDownloadRound : CloudUploadRound))

function publish() {
  $dialog.warning({
    title: "Warning",
    content: `Are you sure you want to ${keyword.value.toLowerCase()} this Application?`,
    positiveText: "Yes, I am sure",
    negativeText: "No, not really",
    onPositiveClick: async () => {
      const data = JSON.parse(JSON.stringify(props.data))
      data.is_published = !data.is_published

      try {
        submitting.value = true

        await http.put(`/api/apps/${props.data.slug}`, data)

        emits("refresh")
        $message.success(`Successfully ${keyword.value.toLowerCase()} this App.`)
      } catch (e: any) {
        $message.error(`Something went wrong... ${e}`)
      } finally {
        submitting.value = false
      }
    },
  })
}
</script>
