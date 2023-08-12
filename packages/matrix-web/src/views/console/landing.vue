<template>
  <div class="container">
    <div class="pt-12 pb-4 px-10">
      <div class="text-2xl font-bold">{{ greetings.period }}, {{ $principal.account?.nickname }}.</div>
      <div class="text-lg">{{ greetings.text }}</div>
    </div>

    <div class="px-8">
      <n-grid x-gap="8" y-gap="8" class="mt-4" item-responsive responsive="screen">
        <n-gi span="24 m:18 l:16">
          <n-card title="Apps">
            <apps />
          </n-card>
        </n-gi>
        <n-gi span="24 m:6 l:8"> </n-gi>
      </n-grid>
    </div>
  </div>
</template>

<script lang="ts" setup>
import Apps from "@/views/console/parts/apps.vue"
import { usePrincipal } from "@/stores/principal"
import { computed } from "vue"

const $principal = usePrincipal()

const greetings = computed(() =>
  new Date().getHours() < 12
    ? { period: "Good Morning", text: "Grab a cup of hot coffee and start your day full of energy!" }
    : new Date().getHours() <= 18 && new Date().getHours() >= 12
    ? { period: "Good Afternoon", text: "It's past three o'clock, what are you doing, the boss won't appreciate you~" }
    : { period: "Good Night", text: "It's time to sleep~ Otherwise the http will explode~" }
)
</script>
