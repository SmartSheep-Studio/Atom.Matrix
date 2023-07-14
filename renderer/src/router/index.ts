import { createRouter, createWebHistory } from "vue-router"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "landing",
      component: () => import("@/views/landing.vue"),
    },
    {
      path: "/store/:app",
      name: "store.info",
      component: () => import("@/views/store/info.vue"),
    },

    {
      path: "/library",
      name: "library",
      component: () => import("@/views/library/landing.vue"),
    },

    {
      path: "/console",
      name: "console",
      component: () => import("@/views/console/landing.vue"),
    },
    {
      path: "/console/:app",
      name: "console.apps",
      component: () => import("@/views/console/app.vue"),
    },
    {
      path: "/console/:app/posts/create",
      name: "console.apps.posts.create",
      component: () => import("@/views/console/actions/create-post.vue"),
    },
    {
      path: "/console/:app/posts/:post/update",
      name: "console.apps.posts.update",
      component: () => import("@/views/console/actions/update-post.vue"),
    },
    {
      path: "/console/:app/releases/create",
      name: "console.apps.releases.create",
      component: () => import("@/views/console/actions/create-release.vue"),
    },
    {
      path: "/console/:app/releases/:release/update",
      name: "console.apps.releases.update",
      component: () => import("@/views/console/actions/update-release.vue"),
    },
    {
      path: "/console/apps/create",
      name: "console.apps.create",
      component: () => import("@/views/console/actions/create-app.vue"),
    },
  ],
})

export default router
