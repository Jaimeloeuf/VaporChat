<script setup lang="ts">
import Logo from './Logo.vue'
import { useWebsocket } from './useWebsocket.ts'

const { websocketConnectionStateString } = useWebsocket()
</script>

<template>
  <div>
    <Logo />
    <p class="text-sm font-light">
      {{ websocketConnectionStateString }}
    </p>
  </div>
</template>
