<template>
  <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="update">
    <n-form-item label="Slug" path="slug">
      <n-input
        placeholder="Used for this Application-Page's link. Can only contain URL-safe characters."
        v-model:value="payload.slug"
      />
    </n-form-item>
    <n-form-item label="URL" path="url">
      <n-input
        placeholder="The homepage of this Application. Can be your studio homepage or source repository, or you can keep this field blank."
        v-model:value="payload.url"
      />
    </n-form-item>
    <n-form-item label="Name" path="name">
      <n-input placeholder="The name of this Application. Accepts anything you want." v-model:value="payload.name" />
    </n-form-item>
    <n-form-item label="Tags" path="tags">
      <n-dynamic-tags v-model:value="payload.tags" />
    </n-form-item>
    <n-form-item label="Description" path="description">
      <n-input
        type="textarea"
        placeholder="A brief description of this Application. Accepts anything you want."
        v-model:value="payload.description"
      />
    </n-form-item>
    <n-form-item label="Details" path="details">
      <v-md-editor v-model="payload.details" height="400px" />
    </n-form-item>

    <n-space size="small">
      <n-button type="primary" attr-type="submit" :loading="submitting">Submit</n-button>
    </n-space>
  </n-form>
</template>

<script lang="ts" setup>
import { parseRedirect } from "@/utils/callback"
import { http } from "@/utils/http"
import { useMessage, type FormRules, type FormInst, useDialog } from "naive-ui"
import { onMounted, reactive, ref } from "vue"
import { useRoute, useRouter } from "vue-router"

const $route = useRoute()
const $router = useRouter()
const $dialog = useDialog()
const $message = useMessage()

const submitting = ref(false)

const props = defineProps<{data: any}>()
const emits = defineEmits(["refresh"])

const form = ref<FormInst | null>(null)
const rules: FormRules = {
  slug: {
    required: true,
    validator: (_, v) => new RegExp(/^[A-Za-z0-9-_]+$/).test(v),
    message: "Only accepts letters, underscores, and numbers",
    trigger: ["blur", "input"],
  },
  name: {
    required: true,
    validator: (_, v) => v.length >= 4,
    message: "Requires at least four characters",
    trigger: ["blur", "input"],
  },
  description: {
    required: true,
    validator: (_, v) => v.length >= 6,
    message: "Requires at least six characters",
    trigger: ["blur", "input"],
  },
  details: {
    required: true,
    validator: (_, v) => v.length >= 6,
    message: "Requires at least six characters",
    trigger: ["blur", "input"],
  },
}

const payload = ref({
  slug: "",
  name: "",
  description: "",
  details: "",
  url: "",
  tags: [],
  is_published: false,
})

function update() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return
    }

    try {
      submitting.value = true

      await http.put(`/api/apps/${payload.value.slug}`, payload.value)

      emits("refresh")
      $message.success("Successfully updated the app")
    } catch (e: any) {
      $message.error(`Something went wrong... ${e}`)
    } finally {
      submitting.value = false
    }
  })
}

onMounted(() => {
  payload.value = props.data
})
</script>
