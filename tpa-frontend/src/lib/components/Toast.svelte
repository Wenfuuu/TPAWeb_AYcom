<script lang="ts">
  import { fade, fly } from "svelte/transition";
  import { createEventDispatcher, onDestroy } from "svelte";

  export let id: string = "";
  export let type: "success" | "error" | "warning" | "info" = "info";
  export let message: string = "";
  export let duration: number = 3000;
  export let position: "top-right" | "top-left" | "bottom-right" | "bottom-left" | "top-center" | "bottom-center" = "top-right";

  const dispatch = createEventDispatcher();
  let timeoutId: ReturnType<typeof setTimeout>;

  function close() {
    clearTimeout(timeoutId);
    dispatch("close", { id });
  }

  $: if (duration > 0) {
    timeoutId = setTimeout(close, duration);
  }

  onDestroy(() => {
    clearTimeout(timeoutId);
  });
</script>

<div
  class="toast {position}"
  in:fly={{ y: position.includes("top") ? -20 : 20, duration: 300 }}
  out:fade={{ duration: 200 }}
>
  <div
    class="toast-content {type}"
    role="alert"
  >
    <div class="toast-icon {type}">
    {#if type === "success"}✓{/if}
    {#if type === "error"}✗{/if}
    {#if type === "warning"}⚠{/if}
    {#if type === "info"}ℹ{/if}
    </div>
    <div class="toast-message">
      <p>{message}</p>
    </div>
    <button
      type="button"
      class="toast-close"
      on:click={close}
      aria-label="Close"
    >
      <span class="sr-only">Close</span>
      <svg width="20" height="20" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path>
      </svg>
    </button>
  </div>
</div>

<style>
  :root {
    --toast-success-bg: #f0fff4;
    --toast-success-border: #48bb78;
    --toast-success-icon: #48bb78;

    --toast-error-bg: #fff5f5;
    --toast-error-border: #f56565;
    --toast-error-icon: #f56565;

    --toast-warning-bg: #fffaf0;
    --toast-warning-border: #ed8936;
    --toast-warning-icon: #ed8936;

    --toast-info-bg: #ebf8ff;
    --toast-info-border: #4299e1;
    --toast-info-icon: #4299e1;
  }

  .toast {
    position: fixed;
    z-index: 50;
    max-width: 350px;
    width: 100%;
  }

  .top-right { top: 1rem; right: 1rem; }
  .top-left { top: 1rem; left: 1rem; }
  .bottom-right { bottom: 1rem; right: 1rem; }
  .bottom-left { bottom: 1rem; left: 1rem; }
  .top-center { top: 1rem; left: 50%; transform: translateX(-50%); }
  .bottom-center { bottom: 1rem; left: 50%; transform: translateX(-50%); }

  .toast-content {
    display: flex;
    align-items: center;
    padding: 12px 16px;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    border-left: 4px solid;
  }

  .toast-content.success {
    background-color: var(--toast-success-bg);
    border-left-color: var(--toast-success-border);
  }

  .toast-content.error {
    background-color: var(--toast-error-bg);
    border-left-color: var(--toast-error-border);
  }

  .toast-content.warning {
    background-color: var(--toast-warning-bg);
    border-left-color: var(--toast-warning-border);
  }

  .toast-content.info {
    background-color: var(--toast-info-bg);
    border-left-color: var(--toast-info-border);
  }

  .toast-icon {
    margin-right: 12px;
    font-size: 1.25rem;
  }

  .toast-icon.success { color: var(--toast-success-icon); }
  .toast-icon.error { color: var(--toast-error-icon); }
  .toast-icon.warning { color: var(--toast-warning-icon); }
  .toast-icon.info { color: var(--toast-info-icon); }

  .toast-message {
    margin-right: 24px;
    color: black;
  }

  .toast-message p {
    font-weight: 500;
    margin: 0;
  }

  .toast-close {
    margin-left: auto;
    background-color: white;
    color: #718096;
    border-radius: 8px;
    padding: 6px;
    height: 32px;
    width: 32px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border: none;
    cursor: pointer;
  }

  .toast-close:hover {
    color: #1a202c;
  }

  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border-width: 0;
  }
</style>
