<script lang="ts">
  import { toastStore } from "../stores/toast";
  import Toast from "./Toast.svelte";
</script>

<div aria-live="polite" aria-atomic="true">
  {#each $toastStore as toast (toast.id)}
    <Toast
      id={toast.id}
      type={toast.type}
      message={toast.message}
      duration={toast.duration}
      position={toast.position}
      on:close={(e) => toastStore.remove(e.detail.id)}
    />
  {/each}
</div>