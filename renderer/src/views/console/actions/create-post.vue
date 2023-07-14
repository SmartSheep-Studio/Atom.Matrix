<template>
  <div class="container">
    <div class="pt-12 pb-4 px-10">
      <div class="text-2xl font-bold">Create a new post</div>
      <div class="text-lg">Don't forgot follow the community guidelines!</div>
    </div>

    <div class="px-10 pt-4">
      <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="create" class="max-w-[800px]">
        <n-form-item label="Slug" path="slug">
          <n-input
            placeholder="Used for this post's link. Can only contain URL-safe characters."
            v-model:value="payload.slug"
          />
        </n-form-item>
        <n-form-item label="Type" path="type">
          <n-select :options="types" v-model:value="payload.type" />
        </n-form-item>
        <n-form-item label="Title" path="title">
          <n-input
            placeholder="The main topic or title of this post. Accepts anything you want."
            v-model:value="payload.title"
          />
        </n-form-item>
        <n-form-item label="Tags" path="tags">
          <n-dynamic-tags v-model:value="payload.tags" />
        </n-form-item>
        <n-form-item label="Content" path="content">
          <v-md-editor v-model="payload.content" height="400px" />
        </n-form-item>

        <n-space size="small">
          <n-button type="primary" attr-type="submit" :loading="submitting">Submit</n-button>
          <n-button @click="$router.back()">Cancel</n-button>
        </n-space>
      </n-form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { parseRedirect } from "@/utils/callback"
import { http } from "@/utils/http"
import { useMessage, type FormRules, type FormInst, useDialog } from "naive-ui"
import { reactive, ref } from "vue"
import { useRoute, useRouter } from "vue-router"

const $route = useRoute()
const $router = useRouter()
const $dialog = useDialog()
const $message = useMessage()

const submitting = ref(false)

const form = ref<FormInst | null>(null)
const rules: FormRules = {
  slug: {
    required: true,
    validator: (_, v) => new RegExp(/^[A-Za-z0-9-_]+$/).test(v),
    message: "Only accepts letters, underscores, and numbers",
    trigger: ["blur", "input"],
  },
  title: {
    required: true,
    validator: (_, v) => v.length >= 4,
    message: "Requires at least four characters",
    trigger: ["blur", "input"],
  },
  content: {
    required: true,
    validator: (_, v) => v.length >= 6,
    message: "Requires at least six characters",
    trigger: ["blur", "input"],
  },
}

const types = [
  { label: "Announcement", value: "announcement" },
  { label: "Blog", value: "blog" },
]

const payload = reactive({
  slug: "",
  type: "announcement",
  title: "",
  content: "",
  tags: [],
  is_published: false,
})

function create() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return
    }

    try {
      submitting.value = true

      await http.post(`/api/apps/${$route.params.app}/posts`, payload)

      $dialog.success({
        title: "Successfully created a post",
        content:
          "Currently your post isn't publish yet, you can publish the post on the console page when you think you're ready.",
        positiveText: "OK",
        onPositiveClick: async () => {
          await $router.push(
            await parseRedirect($route.query, { name: "console.apps", params: { app: $route.params.app } })
          )
        },
      })
    } catch (e: any) {
      $message.error(`Something went wrong... ${e}`)
    } finally {
      submitting.value = false
    }
  })
}
</script>
