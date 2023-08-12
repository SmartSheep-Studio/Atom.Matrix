<template>
  <n-card embedded>
    <div class="text-lg">Destroy App</div>
    <div class="mt-1">
      This operation will destroy all resources that belongs to this application and cannot be undone. Think twice before you press the button.
    </div>

    <div class="mt-3">
      <n-button class="w-full" type="error" :loading="submitting" @click="destroy">
        <template #icon>
          <n-icon :component="DeleteRound" />
        </template>
        Destroy
      </n-button>
    </div>
  </n-card>
</template>

<script lang="ts" setup>
import { http } from "@/utils/http"
import { DeleteRound } from "@vicons/material"
import { useDialog, useMessage } from "naive-ui"
import { ref } from "vue"

const $dialog = useDialog()
const $message = useMessage()

const props = defineProps<{ data: any }>()
const emits = defineEmits(["done"])

const submitting = ref(false)

function destroy() {
  $dialog.warning({
    title: "Warning",
    content: "This operation cannot be undone. Are you sure you want to continue?",
    positiveText: "Yes, I am sure",
    negativeText: "No, not really",
    onPositiveClick: async () => {
      try {
        submitting.value = true

        await http.delete(`/api/apps/${props.data.slug}`)

        emits("done")
        $message.success("Successfully destroyed the App. Redirecting...")
      } catch (e: any) {
        $message.error(`Something went wrong... ${e}`)
      } finally {
        submitting.value = false
      }
    },
  })
}
</script>
