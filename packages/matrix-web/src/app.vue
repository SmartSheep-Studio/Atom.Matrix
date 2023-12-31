<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-dialog-provider>
      <n-message-provider>
        <div class="w-full relative h-[100vh]">
          <n-layout position="absolute">
            <n-layout-header bordered class="px-8 md:px-16 h-full flex items-center">
              <n-menu
                mode="horizontal"
                v-model:value="menuKey"
                :options="menuOptions"
              />
            </n-layout-header>

            <n-layout class="w-full h-[100vh-72px]">
              <data-provider>
                <gatekeeper>
                  <router-view />
                </gatekeeper>
              </data-provider>
            </n-layout>
          </n-layout>
        </div>
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<script lang="ts" setup>
import DataProvider from "@/data-provider.vue";
import Gatekeeper from "@/components/global/gatekeeper.vue";
import { useEndpoint } from "@/stores/endpoint";
import { usePrincipal } from "@/stores/principal";
import { type Component, computed, h, type Ref, ref, watch } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import { type MenuOption, NIcon } from "naive-ui";
import { ExploreRound, LibraryBooksRound, TerminalRound } from "@vicons/material";
import { hasUserPermissions } from "@/utils/gatekeeper";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const $route = useRoute();
const $router = useRouter();
const $endpoint = useEndpoint();
const $principal = usePrincipal();

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

const themeOverrides = {
  common: {
    primaryColor: "#ca4d4dFF",
    primaryColorHover: "#DF5656FF",
    primaryColorPressed: "#C04747FF",
    primaryColorSuppl: "#A84141FF"
  }
};

const menuKey = ref($route.name);
const menuOptions: Ref<MenuOption[]> = computed(() =>
  $principal.isSigned
    ? [
      {
        label: () => h(RouterLink, { to: { name: "landing" } }, { default: () => t("nav.explore") }),
        icon: renderIcon(ExploreRound),
        key: "landing"
      },
      {
        label: () => h(RouterLink, { to: { name: "console" } }, { default: () => t("nav.console") }),
        icon: renderIcon(TerminalRound),
        show: hasUserPermissions("matrix.console.view"),
        key: "console"
      },
      {
        label: () => h(RouterLink, { to: { name: "library" } }, { default: () => t("nav.library") }),
        icon: renderIcon(LibraryBooksRound),
        key: "library"
      }
    ]
    : []
);

watch($route, (v) => {
  menuKey.value = v.name;
});
</script>

<style>
.n-layout-header {
  height: 72px;
}

.n-layout-footer {
  padding: 24px;
}

.w-dialog {
  width: 520px;
}

a {
  color: #329fd2;
}
</style>
