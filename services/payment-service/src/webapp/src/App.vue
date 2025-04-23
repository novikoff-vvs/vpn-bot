<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

window.Telegram.WebApp.expand(); // Раскрываем WebApp на весь экран
window.Telegram.WebApp.enableClosingConfirmation(); // Опционально
window.Telegram.WebApp.MainButton.show(); // Показываем кнопку (если нужно)

const route = useRoute()
const prevDepth = ref(0)
const currentDepth = ref(0)

watch(() => route.path, (newVal, oldVal) => {
  // Определяем "глубину" пути для направления анимации
  prevDepth.value = currentDepth.value
  currentDepth.value = route.meta.depth || 0
})
</script>

<template>
  <div class="container">
      <router-view v-slot="{ Component, route }">
        <transition
            :name="currentDepth > prevDepth ? 'slide-left' : 'slide-right'"
            mode="out-in"
        >
          <component
              :is="Component"
              :key="route.path"
          />
        </transition>
      </router-view>
  </div>
</template>

<style>
/* Анимации перехода */
.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  position: absolute;
  width: 100%;
}

.slide-left-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.slide-left-leave-to {
  transform: translateX(-30%);
  opacity: 0;
}

.slide-right-enter-from {
  transform: translateX(-100%);
  opacity: 0;
}

.slide-right-leave-to {
  transform: translateX(30%);
  opacity: 0;
}

/* Мягкое появление/исчезание для фона */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>