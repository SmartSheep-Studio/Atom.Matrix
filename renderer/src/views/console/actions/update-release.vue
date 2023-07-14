<template>
  <div class="container">
    <div class="pt-12 pb-4 px-10">
      <div class="text-2xl font-bold">Update a exists release</div>
      <div class="text-lg">Don't forgot follow the community guidelines!</div>
    </div>

    <div class="px-10 pt-4">
      <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="update" class="max-w-[800px]">
        <n-form-item label="Slug" path="slug">
          <n-input
            placeholder="Used for this Release's link. Can only contain URL-safe characters."
            v-model:value="payload.slug"
          />
        </n-form-item>
        <n-form-item label="Type" path="type">
          <n-select :options="types" v-model:value="payload.type" />
        </n-form-item>
        <n-form-item label="Name" path="name">
          <n-input placeholder="The name of this Release. Accepts anything you want." v-model:value="payload.name" />
        </n-form-item>
        <n-form-item label="Tags" path="tags">
          <n-dynamic-tags v-model:value="payload.tags" />
        </n-form-item>
        <n-form-item label="Description" path="description">
          <n-input
            type="textarea"
            placeholder="A brief description of this Release. Accepts anything you want."
            v-model:value="payload.description"
          />
        </n-form-item>
        <n-form-item label="Details" path="details">
          <v-md-editor v-model="payload.details" height="400px" />
        </n-form-item>
        <n-form-item v-for="(item, index) in payload.options.assets" :label="`Assets #${index + 1}`">
          <n-input placeholder="Name" v-model:value="item.name" />
          <n-input class="ml-2" placeholder="Download URL" v-model:value="item.url" />
          <n-input class="ml-2" placeholder="Decompressor" v-model:value="item.decompressor" />
          <n-select class="ml-2" :options="platforms" v-model:value="item.platform" />
          <n-button class="ml-2" @click="payload.options.assets.splice(index, 1)">
            <template #icon>
              <n-icon :component="DeleteRound" />
            </template>
          </n-button>
        </n-form-item>
        <n-form-item v-for="(item, index) in payload.options.preprocessing" :label="`Preprocessing #${index + 1}`">
          <n-input placeholder="Name" v-model:value="item.name" />
          <n-input class="ml-2" placeholder="Script" v-model:value="item.script" />
          <n-select class="ml-2" :options="platforms" v-model:value="item.platform" />
          <n-button class="ml-2" @click="payload.options.preprocessing.splice(index, 1)">
            <template #icon>
              <n-icon :component="DeleteRound" />
            </template>
          </n-button>
        </n-form-item>
        <n-form-item v-for="(item, index) in payload.options.run_options" :label="`Run Option #${index + 1}`">
          <n-input placeholder="Name" v-model:value="item.name" />
          <n-input class="ml-2" placeholder="Script" v-model:value="item.script" />
          <n-select class="ml-2" :options="platforms" v-model:value="item.platform" />
          <n-button class="ml-2" @click="payload.options.run_options.splice(index, 1)">
            <template #icon>
              <n-icon :component="DeleteRound" />
            </template>
          </n-button>
        </n-form-item>

        <n-space size="small">
          <n-button @click="payload.options.assets.push({ url: '', name: '', decompressor: '', platform: 'win32' })">
            Add Assets
          </n-button>
          <n-button @click="payload.options.preprocessing.push({ name: '', script: '', platform: 'win32' })">
            Add Preprocessing
          </n-button>
          <n-button @click="payload.options.run_options.push({ name: '', script: '', platform: 'win32' })">
            Add Run Option
          </n-button>
        </n-space>

        <n-form-item class="mt-6" label="Is Published" path="is_published">
          <n-switch v-model:value="payload.is_published" />
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
import { DeleteRound } from "@vicons/material"
import { onMounted, ref } from "vue"
import { useRoute, useRouter } from "vue-router"

const $route = useRoute()
const $router = useRouter()
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

const types = [
  { label: "Minor Update", value: "minor-update" },
  { label: "Major Update", value: "major-update" },
  { label: "Hotfix Update", value: "hotfix-update" },
]

const platforms = [
  { label: "Windows", value: "win32" },
  { label: "MacOS", value: "drawin" },
  { label: "Linux", value: "linux" },
]

const payload = ref<any>({
  slug: "",
  type: "minor-update",
  name: "",
  description: "",
  details: "",
  tags: [],
  options: {
    assets: [],
    preprocessing: [],
    run_options: [],
  },
  is_published: false,
})

async function fetch() {
  try {
    const res = await http.get(`/api/apps/${$route.params.app}/releases/${$route.params.release}`)
    payload.value = res.data
    payload.value.type = res.data.post.type
    payload.value.tags = res.data.post.tags
    payload.value.details = res.data.post.content
  } catch (e: any) {
    $message.error(`Something went wrong... ${e}`)
  }
}

function update() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return
    }

    try {
      submitting.value = true

      await http.put(`/api/apps/${$route.params.app}/releases/${$route.params.release}`, payload.value)

      $message.success("Successfully updated a release")
      await $router.push(
        await parseRedirect($route.query, { name: "console.apps", params: { app: $route.params.app } })
      )
    } catch (e: any) {
      $message.error(`Something went wrong... ${e}`)
    } finally {
      submitting.value = false
    }
  })
}

onMounted(() => {
  fetch()
})
</script>
